package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/proto/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
func ExampleRunClientREST() {
	_ = ExampleClientDevTestsPrintlnREST()
}

func UserGetListREST(httpCli *http.Client, activeselect string) (*model.Users, error) {
	// [GET] - '/users?active=false|true'.
	in := new(model.UserSelect)
	in.Active = activeselect

	reqURL, err := url.Parse(httpDevPath())
	if err != nil {
		return nil, err
	}

	if in.Active != "" {
		query := url.Values{}
		//query.Add("active", activeselect)
		query.Set("active", in.Active)
		reqURL.RawQuery = query.Encode()
	}

	body := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req = httpReqSetHeader(req)

	if httpCli == nil {
		httpCli = NewConnDevHTTP()
	}

	res, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpCli.CloseIdleConnections()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	out := new(model.Users)
	err = json.Unmarshal(resBody, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func UserAddREST(httpCli *http.Client, active bool, name string) (*model.User, error) {
	// [POST] - '/users'
	in := new(model.User)
	in.Active = active
	in.Name = name
	if in.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument: user name is empty")
	}

	reqURL, err := url.Parse(httpDevPath())
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req = httpReqSetHeader(req)

	if httpCli == nil {
		httpCli = NewConnDevHTTP()
	}

	res, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpCli.CloseIdleConnections()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	out := new(model.User)
	err = json.Unmarshal(resBody, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func UserDelREST(httpCli *http.Client, userid int32) (*model.User, error) {
	// [DELETE] - '/users/{user_id}'
	in := new(model.User)
	in.Userid = userid
	reqURL, err := url.Parse(fmt.Sprintf("%s/%d", httpDevPath(), in.Userid))
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req = httpReqSetHeader(req)

	if httpCli == nil {
		httpCli = NewConnDevHTTP()
	}

	res, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpCli.CloseIdleConnections()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	out := new(model.User)
	err = json.Unmarshal(resBody, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func httpDevPath() string {
	// [GET] - '/users?active=false|true'.
	// [DELETE] - '/users/{user_id}'
	// [POST] - '/users'
	if config.ServDev.Secure {
		return "https://" + config.ServDev.RestAddress + "/users"
	} else {
		return "http://" + config.ServDev.RestAddress + "/users"
	}
}

func httpReqSetHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if config.ServDev.Secure {
		req.Header.Set("sys", config.ServDevAuth.Sys)
		req.Header.Set("org", config.ServDevAuth.Org)
		req.Header.Set("login", config.ServDevAuth.Login)
		req.Header.Set("password", config.ServDevAuth.Password)
	}
	return req
}

func NewConnDevHTTP() *http.Client {
	httpCli := http.DefaultClient
	httpCli.Timeout = time.Second * 10
	return httpCli

	// transCfg := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	// }
	// // httpCli := http.DefaultClient
	// httpCli := &http.Client{Transport: transCfg}
	// // httpCli.Transport.RoundTrip()
	// httpCli.Timeout = time.Second * 10
	// return httpCli

	// cert, err := ioutil.ReadFile(config.ServDev.CertFilePath)
	// if err != nil {
	// 	log.Fatal("Couldn't load file", err, config.ServDev.CertFilePath)
	// }
	// certPool := x509.NewCertPool()
	// certPool.AppendCertsFromPEM(cert)
	// transCfg := &http.Transport{
	// 	TLSClientConfig: &tls.Config{RootCAs: certPool}, // ignore expired SSL certificates RootCAs: certPool, InsecureSkipVerify: true
	// }
	// // httpCli := http.DefaultClient
	// httpCli := &http.Client{Transport: transCfg}
	// // httpCli.Transport.RoundTrip()
	// httpCli.Timeout = time.Second * 10
	// return httpCli
}

func ExampleClientDevTestsPrintlnREST() error {
	httpCli := http.DefaultClient
	httpCli.Timeout = time.Duration(10) * time.Second
	defer httpCli.CloseIdleConnections()

	log.Println("start: ExampleClientDevTestsPrintlnREST()")
	tn := time.Now()

	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-mockup")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "true")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "false")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-truncate")))
	// err validation
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, true, ""))) // err empty name
	// add & list
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "true")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, true, "Name-1 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, false, "Name-2 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, false, "Name-3 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, true, "Name-4 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, false, "Name-5 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserAddREST(httpCli, true, "Name-6 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "true")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "false")))
	// del & list
	log.Println(model.DoStringUserErr(UserDelREST(httpCli, 1001)))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserDelREST(httpCli, 1002)))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserDelREST(httpCli, 1003)))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserDelREST(httpCli, 1004)))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	log.Println(model.DoStringUserErr(UserDelREST(httpCli, 1005)))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "")))
	// db commands
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-mockup-100")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-shrink")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-mockup")))
	log.Println(model.DoStringUsersErr(UserGetListREST(httpCli, "dev-db-truncate")))

	log.Println("rest elapsed: ", time.Since(tn))

	return nil
}

package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/implem/cli/auth"
	"github.com/ac-i/user-service/proto/model"
	"github.com/ac-i/user-service/proto/serv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
func ExampleRunClientGRCP() {
	_ = ExampleClientDevTestsPrintlnGRCP()
}

func UserGetListGRPC(grpcConn *grpc.ClientConn, activeselect string) (*model.Users, error) {
	if grpcConn == nil {
		grpcConn, err := NewConnDevGRPC()
		if err != nil {
			return nil, fmt.Errorf("Error grpc.Dial: %s", err)
		}
		defer grpcConn.Close()
	}
	in := new(model.UserSelect)
	in.Active = activeselect
	return serv.NewUserClient(grpcConn).GetList(context.Background(), in)
}

func UserDelGRPC(grpcConn *grpc.ClientConn, userid int32) (*model.User, error) {
	if grpcConn == nil {
		grpcConn, err := NewConnDevGRPC()
		if err != nil {
			return nil, fmt.Errorf("Error grpc.Dial: %s", err)
		}
		defer grpcConn.Close()
	}
	in := new(model.User)
	in.Userid = userid
	return serv.NewUserClient(grpcConn).Del(context.Background(), in)
}

func UserAddGRPC(grpcConn *grpc.ClientConn, active bool, name string) (*model.User, error) {
	if grpcConn == nil {
		grpcConn, err := NewConnDevGRPC()
		if err != nil {
			return nil, fmt.Errorf("Error grpc.Dial: %s", err)
		}
		defer grpcConn.Close()
	}
	in := new(model.User)
	in.Active = active
	in.Name = name
	return serv.NewUserClient(grpcConn).Add(context.Background(), in)
}

func NewConnDevGRPC() (*grpc.ClientConn, error) {
	if config.ServDev.Secure {
		// Create the client TLS credentials
		credsTLS, err := credentials.NewClientTLSFromFile(config.ServDev.CertFilePath, "")
		if err != nil {
			log.Fatalf("could not load tls cert: %s", err)
		}

		// set the login/pass
		authCli := auth.Authentication{
			Sys:      config.ServDevAuth.Sys,
			Org:      config.ServDevAuth.Org,
			Login:    config.ServDevAuth.Login,
			Password: config.ServDevAuth.Password,
		}

		// Initiate a connection with the server WithTransportCredentials & PerRPCCred
		// return grpc.Dial(target, opts...)
		return grpc.Dial(config.ServDev.GrpcAddress,
			grpc.WithTransportCredentials(credsTLS),
			grpc.WithPerRPCCredentials(&authCli))
	} else {
		// Initiate a connection with the server WithInsecure
		return grpc.Dial(config.ServDev.GrpcAddress, grpc.WithInsecure())
	}
}
func ExampleClientDevTestsPrintlnGRCP() error {
	grpcConn, err := NewConnDevGRPC()
	if err != nil {
		return fmt.Errorf("Error grpc.Dial: %s", err)
	}
	defer grpcConn.Close()

	log.Println("start: ExampleClientDevTestsPrintlnGRCP()")
	tn := time.Now()

	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-mockup")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "true")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "false")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-truncate")))
	// err validation
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, true, ""))) // err empty name
	// add & list
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "true")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, true, "Name-1 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, false, "Name-2 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, false, "Name-3 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, true, "Name-4 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, false, "Name-5 Last-F")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserAddGRPC(grpcConn, true, "Name-6 Last-T")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "true")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "false")))
	// del & list
	log.Println(model.DoStringUserErr(UserDelGRPC(grpcConn, 1001)))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserDelGRPC(grpcConn, 1002)))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserDelGRPC(grpcConn, 1003)))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserDelGRPC(grpcConn, 1004)))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	log.Println(model.DoStringUserErr(UserDelGRPC(grpcConn, 1005)))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "")))
	// db commands
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-mockup-100")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-shrink")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-truncate")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-mockup")))
	log.Println(model.DoStringUsersErr(UserGetListGRPC(grpcConn, "dev-db-truncate")))

	log.Println("grpc elapsed: ", time.Since(tn))

	return nil
}

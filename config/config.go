package config

import (
	"log"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

type Serv struct {
	Network         string
	GrpcAddress     string
	RestAddress     string
	Secure          bool
	CertFilePath    string
	KeyFilePath     string
	PresistantStore bool
}

var ServDev = Serv{
	// SetProdMode // now flags in use for now
	Network:     "tcp",
	GrpcAddress: "localhost:8090",
	RestAddress: "localhost:8080",
	// Secure:       true, // comment for false
	// security requires valid key & crt
	KeyFilePath:  "cert/server.key",
	CertFilePath: "cert/server.crt",
	// auth
	// PresistantStore is set in implem/store/userstore.go
	// default dev directory is not fully persistant: os.UserCacheDir()/user-service
	// file name adds suffix representing mode: dev|prod
	// PresistantStore: true, // comment for storage in :memory:
}

// for test and as placeholder for auth
type Authentication struct {
	Sys       string
	Org       string
	Login     string
	Password  string
	AuthAccID string
}

var ServDevAuth = Authentication{
	Sys:       "ss",
	Org:       "oo",
	Login:     "ll",
	Password:  "pp",
	AuthAccID: "aa",
}

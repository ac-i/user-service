package serv

import (
	"log"
	"net"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/implem/handle"
	"github.com/ac-i/user-service/proto/serv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

// main start a gRPC server and waits for connection
func runServerGRCP() (err error) {
	// use default dev config
	cfg := config.ServDev

	// create a listener on TCP port 8090
	netLis, err := net.Listen(cfg.Network, cfg.GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC options for TLS and Auth
	var grpcOpts []grpc.ServerOption
	if cfg.Secure {
		credsTLS, err := credentials.NewServerTLSFromFile(cfg.CertFilePath, cfg.KeyFilePath)
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		grpcOpts = []grpc.ServerOption{grpc.Creds(credsTLS),
			grpc.UnaryInterceptor(unaryInterceptor)}
	}

	// create a gRPC server
	grpcServer := grpc.NewServer(grpcOpts...)
	handleUser := handle.UserServerHandler()
	serv.RegisterUserServer(grpcServer, handleUser)

	// start the server
	if cfg.Secure {
		log.Println("starting secured gRCP server (TLS & auth) on ", netLis.Addr())
	} else {
		log.Println("starting insecure gRCP server (no TLS & no auth) on ", netLis.Addr())
	}
	if err := grpcServer.Serve(netLis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	return nil
}

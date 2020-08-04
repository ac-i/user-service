package serv

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/proto/serv"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func credHeaderMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Sys" || headerName == "Org" || headerName == "Login" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

func runServerRESTGW() (err error) {
	// use default dev config
	cfg := config.ServDev
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var gwMux *runtime.ServeMux

	// REST gw gRPC options for TLS (Auth in client calls)
	var grpcOpts []grpc.DialOption
	if cfg.Secure {
		// configure auth: Incoming Header Matcher and TSL Credentials
		gwMux = runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credHeaderMatcher))
		credsTLS, err := credentials.NewClientTLSFromFile(cfg.CertFilePath, "")
		if err != nil {
			return fmt.Errorf("could not load TLS certificate: %s", err)
		}
		grpcOpts = []grpc.DialOption{grpc.WithTransportCredentials(credsTLS)}
	} else {
		// no TLS in use
		gwMux = runtime.NewServeMux()
		grpcOpts = []grpc.DialOption{grpc.WithInsecure()}
	}

	// Register Service Handler From Endpoint
	err = serv.RegisterUserHandlerFromEndpoint(ctx, gwMux, cfg.GrpcAddress, grpcOpts)
	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}

	// start HTTP/HTTPS REST server
	if cfg.Secure {
		log.Printf("starting secured HTTPS REST server (TLS & auth) on %s", cfg.RestAddress)
		err = http.ListenAndServeTLS(cfg.RestAddress, cfg.CertFilePath, cfg.KeyFilePath, gwMux)
	} else {
		log.Printf("starting insecure HTTP REST server (no TLS & no auth) on %s", cfg.RestAddress)
		err = http.ListenAndServe(cfg.RestAddress, gwMux)
	}
	if err != nil {
		return fmt.Errorf("error http.ListenAndServe : %s", err)
	}
	return nil
}

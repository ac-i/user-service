package serv

import (
	"context"
	"fmt"
	"strings"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/implem/handle"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

// unaryInterceptor calls authenticate Client with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	serv, ok := info.Server.(*handle.UserServer)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	clientID, err := authenticateClient(ctx, serv)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

// authenticateAgent check the client credentials
// TODO - dev test only
func authenticateClient(ctx context.Context, s *handle.UserServer) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// peerC, ok := peer.FromContext(ctx)
		// if ok {
		// 	log.Println("peer.Addr: ", peerC.Addr, ", md: ", md)
		// log.Println("peer.AuthInfo: ", peerC.AuthInfo)
		// }
		clientLogin := strings.Join(md["login"], "")
		// clientPassword := strings.Join(md["password"], "")
		if clientLogin != config.ServDevAuth.Login {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		if strings.Join(md["password"], "") != config.ServDevAuth.Password {
			return "", fmt.Errorf("bad password %s", "clientPassword")
		}
		// clientOrg := strings.Join(md["org"], "")
		// if clientOrg != config.ServDevAuth.Org {
		// 	return "", fmt.Errorf("unknown org %s", clientOrg)
		// }
		// log.Printf("authenticated client: %s, org: %s", clientLogin, clientOrg)
		return config.ServDevAuth.AuthAccID, nil
	}
	return "", fmt.Errorf("missing credentials")
}

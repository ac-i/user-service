package handle

import (
	"context"

	"github.com/ac-i/user-service/implem/store"
	"github.com/ac-i/user-service/proto/model"
	"github.com/ac-i/user-service/proto/serv"
)

var UnimplementedUserServer serv.UnimplementedUserServer

type UserServer struct {
	userStore *store.UserStore
}

func UserServerHandler() *UserServer {
	serv := new(UserServer)
	serv.userStore = new(store.UserStore).Open()
	return serv
}

// DONE GetList - gw get: "/users" ?active=false|true
// PLUS DEV MODE TESTS: active=dev-db-mockup|dev-db-truncate
func (se *UserServer) GetList(ctx context.Context, in *model.UserSelect) (*model.Users, error) {
	_, _ = UnimplementedUserServer.GetList(ctx, in)
	return se.userStore.UserGetList(ctx, in)
}

// DONE Del - gw delete: "/users/{userid}"
func (se *UserServer) Del(ctx context.Context, in *model.User) (*model.User, error) {
	_, _ = UnimplementedUserServer.Del(ctx, in)
	return se.userStore.UserDel(ctx, in)
}

// DONE Add - gw post: "/users"
func (se *UserServer) Add(ctx context.Context, in *model.User) (*model.User, error) {
	_, _ = UnimplementedUserServer.Add(ctx, in)
	return se.userStore.UserAdd(ctx, in)
}

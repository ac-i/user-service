package test

import (
	"testing"

	"github.com/ac-i/user-service/implem/cli"
	"github.com/ac-i/user-service/proto/model"
	"google.golang.org/grpc"
)

// *** All client tests in this file require server running  in background ***
// start server from user-service directory by command: go run cmd/serv/main.go
// ! test commands overwrite data in dev default db on running server

func TestNewConnDevGRPC(t *testing.T) {
	tests := []struct {
		name    string
		want    *grpc.ClientConn
		wantErr bool
	}{
		// DONE: Add test cases.
		{"dev-default-config", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.NewConnDevGRPC()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnDevGRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()
		})
	}
}

func TestUserAddGRPC(t *testing.T) {
	conn, err := cli.NewConnDevGRPC()
	if err != nil {
		t.Errorf("NewConnDevGRPC() error = %v ", err)
		return
	}
	defer conn.Close()
	type args struct {
		grpcConn *grpc.ClientConn
		active   bool
		name     string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "Name-1 Last-T", args: args{grpcConn: conn, active: true, name: "Name-1 Last-T"}, want: &model.User{Active: true, Name: "Name-1 Last-T"}, wantErr: false},
		{name: "Name-2 Last-F", args: args{grpcConn: conn, active: false, name: "Name-2 Last-F"}, want: &model.User{Active: false, Name: "Name-2 Last-F"}, wantErr: false},
		{name: "Name-3 Last-F", args: args{grpcConn: conn, active: false, name: "Name-3 Last-F"}, want: &model.User{Active: false, Name: "Name-3 Last-F"}, wantErr: false},
		{name: "Name-4 Last-T", args: args{grpcConn: conn, active: true, name: "Name-4 Last-T"}, want: &model.User{Active: true, Name: "Name-4 Last-T"}, wantErr: false},
		{name: "Name-5 Last-F", args: args{grpcConn: conn, active: false, name: "Name-5 Last-F"}, want: &model.User{Active: false, Name: "Name-5 Last-F"}, wantErr: false},
		{name: "Name-6 Last-T", args: args{grpcConn: conn, active: true, name: "Name-6 Last-T"}, want: &model.User{Active: true, Name: "Name-6 Last-T"}, wantErr: false},
		// err
		{name: "empty name", args: args{grpcConn: conn, active: true, name: ""}, want: &model.User{Active: true, Name: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.UserAddGRPC(tt.args.grpcConn, tt.args.active, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAddGRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Active != tt.want.Active {
					t.Errorf("UserAddGRPC() = %v, want %v", got.Active, tt.want.Active)
				}
				if got.Name != tt.want.Name {
					t.Errorf("UserAddGRPC() = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}
}

func TestUserDelGRPC(t *testing.T) {
	conn, err := cli.NewConnDevGRPC()
	if err != nil {
		t.Errorf("NewConnDevGRPC() error = %v ", err)
		return
	}
	defer conn.Close()
	// truncate db
	_, _ = cli.UserGetListGRPC(conn, "dev-db-truncate")
	// add 10 new by mockup (uderid starts from 1001)
	_, _ = cli.UserGetListGRPC(conn, "dev-db-mockup")

	type args struct {
		grpcConn *grpc.ClientConn
		userid   int32
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "del-0", args: args{grpcConn: conn, userid: 0}, want: &model.User{Userid: 0}, wantErr: true},
		{name: "del-1", args: args{grpcConn: conn, userid: 1}, want: &model.User{Userid: 0}, wantErr: false},
		{name: "del-1001", args: args{grpcConn: conn, userid: 1001}, want: &model.User{Userid: 1001}, wantErr: false},
		{name: "del-1001-2", args: args{grpcConn: conn, userid: 1001}, want: &model.User{Userid: 0}, wantErr: false},
		{name: "del-1002", args: args{grpcConn: conn, userid: 1002}, want: &model.User{Userid: 1002}, wantErr: false},
		{name: "del-1004", args: args{grpcConn: conn, userid: 1004}, want: &model.User{Userid: 1004}, wantErr: false},
		{name: "del-1006", args: args{grpcConn: conn, userid: 1006}, want: &model.User{Userid: 1006}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.UserDelGRPC(tt.args.grpcConn, tt.args.userid)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDelGRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Userid != tt.want.Userid {
					t.Errorf("UserAddGRPC() = %v, want %v", got.Userid, tt.want.Userid)
				}
			}
		})
	}
}

func TestUserGetListGRPC(t *testing.T) {
	conn, err := cli.NewConnDevGRPC()
	if err != nil {
		t.Errorf("NewConnDevGRPC() error = %v ", err)
		return
	}
	defer conn.Close()

	type args struct {
		grpcConn     *grpc.ClientConn
		activeselect string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Users
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "all", args: args{grpcConn: conn, activeselect: ""}, want: nil, wantErr: false},
		{name: "dev-db-truncate", args: args{grpcConn: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		{name: "dev-db-mockup", args: args{grpcConn: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
		{name: "all", args: args{grpcConn: conn, activeselect: ""}, want: nil, wantErr: false},
		{name: "active=true", args: args{grpcConn: conn, activeselect: "true"}, want: nil, wantErr: false},
		{name: "active=false", args: args{grpcConn: conn, activeselect: "false"}, want: nil, wantErr: false},
		// err
		{name: "err", args: args{grpcConn: conn, activeselect: "err"}, want: nil, wantErr: true},
		// cleanup db
		{name: "dev-db-shrink", args: args{grpcConn: conn, activeselect: "dev-db-shrink"}, want: nil, wantErr: false},
		{name: "dev-db-truncate", args: args{grpcConn: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		// populate db - can be useful to keep results (works only on empty db /easy by truncate)
		// {name: "dev-db-truncate", args: args{grpcConn: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		// {name: "dev-db-mockup", args: args{grpcConn: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
		// {name: "dev-db-mockup-100", args: args{grpcConn: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.UserGetListGRPC(tt.args.grpcConn, tt.args.activeselect)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGetListGRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

package test

import (
	"net/http"
	"testing"

	"github.com/ac-i/user-service/implem/cli"
	"github.com/ac-i/user-service/proto/model"
)

// *** All client tests in this file require server running  in background ***
// start server from user-service directory by command: go run cmd/serv/main.go
// ! test commands overwrite data in dev default db on running server

func TestUserAddREST(t *testing.T) {
	conn := cli.NewConnDevHTTP()
	defer conn.CloseIdleConnections()
	type args struct {
		httpCli *http.Client
		active  bool
		name    string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "Name-1 Last-T", args: args{httpCli: conn, active: true, name: "Name-1 Last-T"}, want: &model.User{Active: true, Name: "Name-1 Last-T"}, wantErr: false},
		{name: "Name-2 Last-F", args: args{httpCli: conn, active: false, name: "Name-2 Last-F"}, want: &model.User{Active: false, Name: "Name-2 Last-F"}, wantErr: false},
		{name: "Name-3 Last-F", args: args{httpCli: conn, active: false, name: "Name-3 Last-F"}, want: &model.User{Active: false, Name: "Name-3 Last-F"}, wantErr: false},
		{name: "Name-4 Last-T", args: args{httpCli: conn, active: true, name: "Name-4 Last-T"}, want: &model.User{Active: true, Name: "Name-4 Last-T"}, wantErr: false},
		{name: "Name-5 Last-F", args: args{httpCli: conn, active: false, name: "Name-5 Last-F"}, want: &model.User{Active: false, Name: "Name-5 Last-F"}, wantErr: false},
		{name: "Name-6 Last-T", args: args{httpCli: conn, active: true, name: "Name-6 Last-T"}, want: &model.User{Active: true, Name: "Name-6 Last-T"}, wantErr: false},
		// err
		{name: "empty name", args: args{httpCli: conn, active: true, name: ""}, want: &model.User{Active: true, Name: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.UserAddREST(tt.args.httpCli, tt.args.active, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAddREST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Active != tt.want.Active {
					t.Errorf("UserAddREST() = %v, want %v", got.Active, tt.want.Active)
				}
				if got.Name != tt.want.Name {
					t.Errorf("UserAddREST() = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}
}

func TestUserDelREST(t *testing.T) {
	conn := cli.NewConnDevHTTP()
	defer conn.CloseIdleConnections()
	// truncate db
	_, _ = cli.UserGetListREST(conn, "dev-db-truncate")
	// add 10 new by mockup (uderid starts from 1001)
	_, _ = cli.UserGetListREST(conn, "dev-db-mockup")

	type args struct {
		httpCli *http.Client
		userid  int32
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "del-0", args: args{httpCli: conn, userid: 0}, want: &model.User{Userid: 0}, wantErr: false}, //? err rest on empty
		{name: "del-1", args: args{httpCli: conn, userid: 1}, want: &model.User{Userid: 0}, wantErr: false},
		{name: "del-1001", args: args{httpCli: conn, userid: 1001}, want: &model.User{Userid: 1001}, wantErr: false},
		{name: "del-1001-2", args: args{httpCli: conn, userid: 1001}, want: &model.User{Userid: 0}, wantErr: false},
		{name: "del-1002", args: args{httpCli: conn, userid: 1002}, want: &model.User{Userid: 1002}, wantErr: false},
		{name: "del-1004", args: args{httpCli: conn, userid: 1004}, want: &model.User{Userid: 1004}, wantErr: false},
		{name: "del-1006", args: args{httpCli: conn, userid: 1006}, want: &model.User{Userid: 1006}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.UserDelREST(tt.args.httpCli, tt.args.userid)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDelREST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Userid != tt.want.Userid {
					t.Errorf("UserAddREST() = %v, want %v", got.Userid, tt.want.Userid)
				}
			}
		})
	}
}

func TestUserGetListREST(t *testing.T) {
	conn := cli.NewConnDevHTTP()
	defer conn.CloseIdleConnections()

	type args struct {
		httpCli      *http.Client
		activeselect string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Users
		wantErr bool
	}{
		// DONE: Add test cases.
		{name: "all", args: args{httpCli: conn, activeselect: ""}, want: nil, wantErr: false},
		{name: "dev-db-truncate", args: args{httpCli: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		{name: "dev-db-mockup", args: args{httpCli: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
		{name: "all", args: args{httpCli: conn, activeselect: ""}, want: nil, wantErr: false},
		{name: "active=true", args: args{httpCli: conn, activeselect: "true"}, want: nil, wantErr: false},
		{name: "active=false", args: args{httpCli: conn, activeselect: "false"}, want: nil, wantErr: false},
		// err
		{name: "err", args: args{httpCli: conn, activeselect: "err"}, want: nil, wantErr: false}, // ? err on rest
		// cleanup db
		{name: "dev-db-shrink", args: args{httpCli: conn, activeselect: "dev-db-shrink"}, want: nil, wantErr: false},
		{name: "dev-db-truncate", args: args{httpCli: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		// populate db - can be useful to keep results (works only on empty db /easy by truncate)
		// {name: "dev-db-truncate", args: args{httpCli: conn, activeselect: "dev-db-truncate"}, want: nil, wantErr: false},
		// {name: "dev-db-mockup", args: args{httpCli: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
		// {name: "dev-db-mockup-100", args: args{httpCli: conn, activeselect: "dev-db-mockup"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cli.UserGetListREST(tt.args.httpCli, tt.args.activeselect)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGetListREST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

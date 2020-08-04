# Test info

## Server gRCP & REST gateway
	server fires separate go routines for gRCP & REST gateway
	start server in separate terminal:
	change directory to user-service
	run server by: go run .
	it will fire serv.RunServerMAIN()
```sh
go run .
# subdirectory cmd contains main packages for server and clients
# you can start the server like:
# go run cmd/serv/*
# go run cmd/serv/main.go
```
```go
// go run . is fires:
serv.RunServerMAIN()
```

### Client gRCP & REST gateway
	To run client examples in go just enter: go run cmd/cli/*
	It will fire: cli.ExampleRunClientGRCP() & cli.ExampleRunClientREST()
	which prints on tne screen results of calling sequences of the user service methods 
```sh
go run cmd/cli/*
# go run cmd/cli/main.go
```
```go
cli.ExampleRunClientGRCP()
cli.ExampleRunClientREST()
```

### Client REST by curl - simple test
```sh
curl -X GET "http://localhost:8080/users"
curl -X GET "http://localhost:8080/users?active=true"
curl -X GET "http://localhost:8080/users?active=false"
curl -X POST -d '{"active":true,"name":"Name-1 Last-T"}' "http://localhost:8080/users?active=true"
curl -X POST -d '{"active":false,"name":"Name-2 Last-F"}' "http://localhost:8080/users?active=true"
curl -X GET "http://localhost:8080/users"
curl -X GET "http://localhost:8080/users?active=true"
curl -X GET "http://localhost:8080/users?active=false"
curl -X GET "http://localhost:8080/users?active="
curl -X GET "http://localhost:8080/users?active=1"
curl -X GET "http://localhost:8080/users?active=0"
curl -X DELETE "http://localhost:8080/users/1001"
curl -X GET "http://localhost:8080/users"
curl -X GET "http://localhost:8080/users?active=true"
curl -X GET "http://localhost:8080/users?active=false"
curl -X DELETE "http://localhost:8080/users/1002"
curl -X GET "http://localhost:8080/users"
curl -X GET "http://localhost:8080/users?active=true"
curl -X GET "http://localhost:8080/users?active=false"
# to clear storage db use truncate dev hack: 
# curl -X GET "http://localhost:8080/users?active=dev-db-truncate"
```
### Client REST by curl - advanced playground with 
```sh
# edit comment/uncomment and fire tests/curl-cli-rest-gw.sh
tests/curl-cli-rest-gw.sh
```
```sh
# ********* gRCP-gw-REST client playground *********
# 
NODE="http://localhost:8080"
# NODE="http://192.168.1.53:8080"
# 
HEAD=""
# 
# for secure TLS (use in pair with secure HEAD below)
# NODE="https://localhost:8080"
#
# for secure TLS and auth
# HEAD=" --cacert "cert/server.crt" -H "sys:ss" -H "org:oo" -H "login:ll" -H "password:pp" "
# for secure TLS (no cert verification) and auth
# HEAD=" -k -H "sys:ss" -H "org:oo" -H "login:ll" -H "password:pp" "
# 
START=$(date +"%s%N")
clear
echo "------START------ ${START} >>> ${NODE} >>> "
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo && echo "--ADD--&--LIST--"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-1 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":false,"name":"Name-2 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"name":"Name-3 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-4 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":false,"name":"Name-5 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-6 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo && echo "--DELETE--&--LIST--"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1001"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1001"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1002"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1002"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1003"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1004"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1005"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo && echo "--DB--"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo "--mck--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-mockup"
echo && echo "--shr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-shrink"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo "--mck--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-mockup-100"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
# 
END=$(date +"%s%N")
TNS=$(((${END}-${START})/1))
TUS=$(((${END}-${START})/1000))
TMS=$(((${END}-${START})/1000000))
TS=$(((${END}-${START})/1000000000))
# 
echo && echo && echo "------END------ ${END} >>> ${NODE} >>> (${TS} s, ${TMS} ms, ${TUS} us, ${TNS} ns)"
# 
# ********* playground end *********
# 
# 
# -------- notes
# 
# REST GW insecure (no TLS handshake, no user/auth interceptor)
# memory
# ------END------ 1596370176324721217 >>> localhost >>> (0 s, 218 ms, 218584 us, 218584633 ns)
# ------END------ 1596370120923700532 >>> 192.168.1.53 >>> (0 s, 263 ms, 263097 us, 263097088 ns)
# persist
# ------END------ 1596370280772783537 >>> localhost >>> (2 s, 2516 ms, 2516014 us, 2516014986 ns)
# ------END------ 1596370405797706301 >>> 192.168.1.53 >>> (2 s, 2645 ms, 2645307 us, 2645307007 ns)
# 
# REST GW secure (with TLS handshake, with user/auth interceptor per call)
# memory
# ------END------ 1596372487554309238 >>> localhost >>> (0 s, 225 ms, 225900 us, 225900406 ns)
# persist
# ------END------ 1596372600525864317 >>> localhost >>> (2 s, 2514 ms, 2514646 us, 2514646711 ns)
```
# The Initial Task

	create a microservice in Go. Simple API service, let's call it "user-service". It would serve 3 endpoints:

	[GET] - '/users?active=false|true'.
	Fetches all users from the database and return them in 'application/json' format.
	It would take one optional 'query parameter' called 'active'(boolean type). If provided return only users with 'active' field set to given value true or false, if not - return all users.

	[DELETE] - '/users/{user_id}'
	Deletes user with 'user_id' provided as 'path parameter'

	[POST] - '/users'
	Creates new user. request body: {name string, active boolean}, and returns full user data{name, user_id, active}
	UserID should be generated either by you or by database.

	USER structure:
	UserID (type up to you)
	Active bool
	Name string

	For dependency managment use Go Modules.
	Database you are gonna use is up to you. Might be MongoDB, PostgreSQL or what you simply feel strong with.
	Library for handling routing is also up to you. Might be chi-mux, gin-gonic or simple http handlers in native 'http' package.
	Please, create some public repository on your github and push the implementation there. After cloning your solution I would like to be able to do two things:
	*******
	Run the tests with go test command
	Run the application with go run command
	Good luck!

# Implementation notes

	1. The USER data model and service Methods  defined with Protocol Buffers in separate model and serv files for better clarity
	2. Automatically generated gRPC and REST-gateway stubs, containing strongly typed model, service methods and REST endpoints
	3. Data storage /persistance (or in memory dev default) uses embedded BuntDB
	4. Written server handlers for user service methods
	5. Written gRPC and REST clients
	6. Other
		a) added mockup and truncate commands for more fun with test :)
		b) added unrequested 'validation' of user name to show/test error messages


	Service User Methods:

	GetList    - supports http/rest: [GET] - '/users?active=false|true'.
	fetches all users or filters by query parameter ?active=false|true
	overcoming false vs empty of a boolean type variant by using string
	active=false|true|'' - to filter false or true or all
	PLUS DEV MODE TESTS: active=dev-db-mockup|dev-db-truncate
	GetList(context.Context, *model.UserSelect) (*model.Users, error)

	Del        - supports http/rest: [DELETE] - '/users/{user_id}'
	deletes record by userid
	returns deleted record or error message
	Del(context.Context, *model.User) (*model.User, error)

	Add        - supports http/rest: [POST] - '/users'
	adds user record, require name!='', field active=true is optional
	userid is calculated by db (at save transaction as max value + 1)
	Add(context.Context, *model.User) (*model.User, error)
```proto
// User data model

message User {
  int32 userid = 1;
  bool active = 2;
  string name = 3;
}

message UserSelect {
  string active = 4;
}

message Users {
  repeated User users = 1;
}

// Service User methods

service User {
  rpc GetList(model.UserSelect) returns (model.Users) {
    option (google.api.http) = {
        get: "/users"
    };
  }
  rpc Del(model.User) returns (model.User) {
    option (google.api.http) = {
        delete: "/users/{userid}"
    };
  }
  rpc Add(model.User) returns (model.User) {
    option (google.api.http) = {
        post: "/users"
        body: "*"
    };
  }
}
```
```json
// configuration VS Code (settings.json) to generate pb on the fly
{
  "protoc": {
    "path": "${env.HOME}/pb/bin/protoc",
    "compile_on_save": true,
    "options": [
      "--proto_path=.",
      "--proto_path=${env.HOME}/gt/ac-i/",
      "--proto_path=${env.HOME}/pb/include",
      "--go_out=plugins=grpc:${env.HOME}/gt/ac-i/",
      "--grpc-gateway_out=logtostderr=true:${env.HOME}/gt/ac-i/",
    //   "--swagger_out=logtostderr=true:."
    ]
  }
}

```

# Have fun :)
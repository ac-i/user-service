package main

import (
	serv "github.com/ac-i/user-service/implem/serv"
)

// main start a gRPC server and waits for connection
func main() {
	serv.RunServerMAIN()
}

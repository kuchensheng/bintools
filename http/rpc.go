package http

import (
	"fmt"
	"net/rpc"
)

var rpcServer = false

//Register publishes the receiver's methods in the DefaultServer
func (e *engine) Register(target any) {
	rpc.Register(target)
	rpcServer = true
}

func (e *engine) RunRpc() {
	if rpcServer {
		fmt.Println("start server with rpc support")
		rpc.HandleHTTP()
	}
}

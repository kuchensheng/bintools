package http

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

//Register publishes the receiver's methods in the DefaultServer
func (e *engine) Register(name string, target any) {
	rpc.RegisterName(name, target)
	e.RpcServer = true
}

func (e *engine) RunRpc(port int) {
	if e.RpcServer {
		fmt.Printf("starting server with rpc support,port:%d...\n", port)
		if l, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
			log.Fatal("can not start rpc server,err is ", err)
		} else {
			go rpc.Accept(l)
			rpc.HandleHTTP()
		}
	}
}

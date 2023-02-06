package http

import (
	"log"
	"net/rpc"
	"testing"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

// Some of Arith's methods have value args, some have pointer args. That's deliberate.

func (t *Arith) Add(args Args, reply *Reply) error {
	reply.C = args.A + args.B
	println("reply.C = ", reply.C)
	return nil
}

func TestRegister(t *testing.T) {
	e := Default()
	e.Register("Arith", new(Arith))
	e.Run(8080)
}

func TestRpcClient(t *testing.T) {
	args := &Args{7, 8}
	reply := new(Reply)

	if c, err := rpc.Dial("tcp", "127.0.0.1:8081"); err != nil {
		log.Fatalf("dialing:%v", err)
	} else if err = c.Call("Arith.Add", args, reply); err != nil {
		log.Printf("error : %v", err)
	} else {
		c.Close()
		log.Printf("侬好")
	}
}

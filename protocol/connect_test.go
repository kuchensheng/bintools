package protocol

import (
	"fmt"
	"github.com/kuchensheng/bintools/logger"
	"github.com/kuchensheng/protocol/metadata"
	"net"
	"testing"
)

func TestConn_Read(t *testing.T) {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	log = logger.GlobalLogger
	for true {
		c, e := listen.Accept()
		if e != nil {
			log.Error("无法获取conn信息,%v", e)
			continue
		}
		conn := Conn{c, metadata.Metadata{}, nil}
		go func() {
			read, err2 := conn.Read()
			if err2 != nil {
				//log.Error("conn read error,%s", err2.Error())
				fmt.Errorf("%v", err2)
				return
			}
			log.Info("%s", read)
		}()
	}
}

func TestConn_Write(t *testing.T) {
	client, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	t.Log("client connect")
	defer client.Close()

	msg := "我是库陈胜"
	_, err = client.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	t.Log("client disconnect")
}

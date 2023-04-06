package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/kuchensheng/bintools/logger"
	"github.com/kuchensheng/protocol/metadata"
	"math"
	"net"
)

var log = logger.GlobalLogger

const (
	ConstHeader       = "Isc-Header"
	ConstHeaderLength = 7
	ConstData         = "Isc-Data"
	MagicFirst        = 0x01
	MagicSecond       = 0x02
	MinLen            = 4
	MaxLen            = math.MaxUint16
)

var magic = []byte{MagicFirst, MagicSecond}

type Conn struct {
	conn net.Conn
	meta metadata.Metadata
	scan *bufio.Scanner
}

func (c *Conn) Read() (b []byte, err error) {
	fmt.Sprintf("%s", "读取连接信息")
	//log.Debug("读取链接的信息,meta={%v}", c.meta)
	if c == nil {
		return nil, errors.New("conn is nil")
	}
	if c.scan == nil {
		c.RegisterSplitFunc()
	}
	if c.scan.Scan() {
		b = c.scan.Bytes()
	} else {
		err = c.conn.Close()
		if err != nil {
			return nil, err
		}
		err = errors.New("conn scan error")
	}
	return
}

func (c *Conn) Write(data []byte) (int, error) {
	if len(data) > MaxLen {
		return 0, errors.New("data too long")
	}

	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, magic)
	if err != nil {
		return 0, err
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(len(data)+MinLen))
	if err != nil {
		return 0, err
	}
	err = binary.Write(&buf, binary.BigEndian, data)
	if err != nil {
		return 0, err
	}
	return c.conn.Write(buf.Bytes())
}

func (c *Conn) RegisterSplitFunc() {
	//todo 做一个scan的split操作
	scanner := bufio.NewScanner(c.conn)
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//判断长度,魔数
		if len(data) < MinLen || data[0] != MagicFirst || data[1] != MagicSecond {
			return
		}
		var l uint16
		//获取header
		err = binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &l)
		if err != nil {
			log.Error("binary.Read Error:%v", err)
			return
		}
		//通过长度读取数据,advance为读取的长度，包括头和数据，data是读取的数据
		if int(l) <= len(data) {
			advance, token, err = int(l), data[:int(l)], nil
		}
		if atEOF {
			err = errors.New("EOF")
		}
		return
	}
	scanner.Split(splitFunc)
	c.scan = scanner
}

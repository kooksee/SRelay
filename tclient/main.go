package tclient

import (
	"bufio"
	"bytes"
	"time"

	kts "SRelay/types"
	knet "SRelay/utils/net"
)

// Run app run
func Run() {

	if c, err := knet.ConnectTcpServer(":8081"); err != nil {
		panic(err.Error())
	} else {
		go handler(c)
	}
}

func handler(c knet.Conn) {
	c.SetReadDeadline(time.Now().Add(connReadTimeout))
	c.SetReadDeadline(time.Time{})
	read := bufio.NewReader(c)

	msg := &kts.KMsg{
		Account: "123456",
		Event:   "account",
		Token:   "123456",
	}

	// 注册客户端信息
	c.Write(msg.Dumps())

	for {

		message, err := read.ReadBytes('\n')
		if err != nil {
			log.Info("tcp error ", err.Error())
			break
		}
		message = bytes.Trim(message, "\n")
		log.Info(string(message))
	}
}

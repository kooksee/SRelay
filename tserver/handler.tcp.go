package tserver

import (
	"bufio"
	"bytes"
	"time"

	knet "SRelay/utils/net"

	kts "SRelay/types"

	"github.com/json-iterator/go"
)

func TcpHandleListener(l *knet.TcpListener) {

	var message []byte

	for {
		c, err := l.Accept()
		if err != nil {
			log.Error("Listener for incoming connections from client closed")
			return
		}

		// Start a new goroutine for dealing connections.
		go func(conn knet.Conn) {
			conn.SetReadDeadline(time.Now().Add(connReadTimeout))
			conn.SetReadDeadline(time.Time{})
			read := bufio.NewReader(conn)
			for {

				message, err = read.ReadBytes('\n')
				if err != nil {
					log.Info("tcp error ", err.Error())
					break
				}
				message = bytes.Trim(message, "\n")

				log.Info("tcp msg ", string(message))

				// 解析请求数据
				msg := &kts.KMsg{}
				if err := jsoniter.Unmarshal(message, msg); err != nil {
					conn.Write(kts.ResultError(err.Error()))
					continue
				}
			}
		}(c)
	}
}

package tserver

import (
	"strings"
	"time"

	knet "SRelay/utils/net"
)

type KcpManager struct {
	clients map[string]knet.Conn
	l       *knet.KcpListener

	host string
	port int
}

func NewKcpManager(host string, port int) *KcpManager {
	return &KcpManager{host: host, port: port}
}

func (km *KcpManager) Listen() (err error) {
	km.l, err = knet.ListenKcp(km.host, km.port)
	return
}

func (km *KcpManager) Start() {
	go func() {
		for {
			c, err := km.l.Accept()
			c.SetReadDeadline(time.Now().Add(connReadTimeout))
			c.SetReadDeadline(time.Time{})

			addr := strings.Split(c.LocalAddr().String(), ":")
			if err != nil {
				delete(km.clients, addr[1])
				continue
			}
			km.clients[addr[1]] = c
		}
	}()
}

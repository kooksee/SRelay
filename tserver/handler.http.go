package tserver

import (
	"time"

	knet "SRelay/utils/net"
)

type HttpManager struct {
	l *knet.KcpListener

	host string
	port int
}

func NewHttpManager(host string, port int) *HttpManager {
	return &HttpManager{host: host, port: port}
}

func (km *HttpManager) Listen() (err error) {
	km.l, err = knet.ListenKcp(km.host, km.port)
	return
}

func (km *HttpManager) Start() {

	go func() {
		for {
			c, err := km.l.Accept()
			if err != nil {

			}
			c.SetReadDeadline(time.Now().Add(connReadTimeout))

			
		}
	}()

}

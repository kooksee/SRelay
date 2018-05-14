package tserver

import (
	"fmt"

	knet "SRelay/utils/net"
)

// Run app run
func Run() {

	clients = make(map[int][]knet.Conn)

	if listener, err := knet.ListenTcp(":8082"); err != nil {
		panic(fmt.Sprintf("Create server listener error, %v", err.Error()))
	} else {
		go TcpHandleListener(listener)
	}

	if l, err := knet.ListenUDP("8083"); err != nil {
		panic(fmt.Sprintf("Create server listener error, %v", err.Error()))
	} else {
		go UdpHandleListener(l)
	}

	l, err := knet.ListenKcp("0.0.0.0", 1234)
	if err != nil {

	}
}

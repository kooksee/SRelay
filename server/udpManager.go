package server

import (
	"bufio"
	"fmt"

	"github.com/kooksee/srelay/utils"
	knet "github.com/kooksee/srelay/utils/net"

	kts "github.com/kooksee/srelay/types"
)

type UdpServerManager struct {
	lis  *knet.UdpListener
	ks   *KcpServer
	port int
}

func (u *UdpServerManager) CreateUdp(port int) error {

	l, err := knet.ListenUDP(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	u.lis = l

	go u.onHandleListen(l)
	return nil
}

func (u *UdpServerManager) onHandleConn(conn knet.Conn) {
	for {
		//conn.RemoteAddr().String()

		r := bufio.NewReader(conn)
		message, err := r.ReadBytes(kts.Delim)
		if err != nil {
			// 节点退出
			logger.Error(err.Error())
			break
		}

		message = utils.BytesTrimSpace(message)

		// 从message中解析出一个UUID出来，然后缓存这个UUID
		//message

		// 数据转发给客户端
		c := u.ks.clients[u.port]
		c.Write(message)

	}
}
func (u *UdpServerManager) onHandleListen(l *knet.UdpListener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		go u.onHandleConn(conn)
	}
}

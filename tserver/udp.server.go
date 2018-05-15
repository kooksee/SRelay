package tserver

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/kataras/iris/core/errors"
	knet "github.com/kooksee/srelay/utils/net"

	kts "github.com/kooksee/srelay/types"
)

type UdpServerManager struct {
	MaxPort int
	MinPort int
	umap    map[int]*knet.UdpListener
}

func (u *UdpServerManager) CreateUdp() (*net.UDPAddr, error) {
	if len(u.umap) > u.MaxPort-u.MinPort {
		return nil, errors.New("端口数量操作设置的最大限制")
	}

	for i := u.MinPort; i < u.MaxPort; i++ {
		if _, ok := u.umap[i]; ok {
			continue
		}
		l, err := knet.ListenUDP(fmt.Sprintf("0.0.0.0:%d", i))
		if err != nil {
			return nil, err
		}
		u.umap[i] = l
		go u.onHandleListen(u.umap[i])
		return &net.UDPAddr{Port: i}, nil
	}
	return nil, errors.New("未知错误")
}

func (u *UdpServerManager) onHandleConn(conn knet.Conn) {
	r := bufio.NewReader(conn)
	for {
		message, err := r.ReadBytes(kts.Delim)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		message = bytes.Trim(message, string(kts.Delim))

		// 解析请求数据
		msg := &kts.KMsg{}
		if err := json.Unmarshal(message, msg); err != nil {
			logger.Error(err.Error())
			continue
		}
		ksInstance.Send(msg)
	}
}
func (u *UdpServerManager) onHandleListen(l *knet.UdpListener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		go u.onHandleConn(conn)
	}
}

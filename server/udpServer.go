package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kooksee/srelay/types"
)

type UdpServer struct {
	l *net.UDPAddr
}

func (u *UdpServer) Listen(port int64) error {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	readConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	go u.onHandleConn(readConn)
	return nil
}

func (u *UdpServer) onHandleConn(conn *net.UDPConn) {
	kb := types.NewKBuffer()
	for {
		buf := make([]byte, 1024*16)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("onHandleConn1", "err", err.Error())

			if err == io.EOF {
				break
			}

			time.Sleep(time.Second * 2)
			continue
		}

		if cfg.Debug {
			logger.Debug(string(bytes.TrimSpace(buf[:n])), "msg", "udp message", "addr", addr.String())
		}

		messages := kb.Next(buf[:n])
		if messages == nil {
			continue
		}

		for _, m := range messages {
			if m == nil || len(m) == 0 {
				continue
			}

			//	发送数据给客户端
			c, err := types.DecodeClient(m)
			if err != nil {
				logger.Error(string(types.ErrJsonParse(err)), "data", string(m))

				//	数据解析失败，发送给udp客户端
				if _, err := conn.WriteToUDP(types.ErrJsonParse(err), addr); err != nil {
					logger.Error("onHandleConn2", "err", err.Error())
				}
				continue
			}

			// 得到id
			con, b := clientsCache.Get(c.TID)
			if !b {
				if _, err := conn.WriteToUDP(types.ErrPeerNotFound(errors.New(fmt.Sprintf("peer %s is nonexistent", c.TID))), addr); err != nil {
					logger.Error("onHandleConn3", "err", err.Error())
				}
				continue
			}

			if _, err := con.(*net.TCPConn).Write(append(m, "\n"...)); err != nil {
				//	写入后端数据失败
				if _, err := conn.WriteToUDP(types.ErrPeerWrite(errors.New(fmt.Sprintf("peer %s write %s error", c.TID, m[1:]))), addr); err != nil {
					logger.Error("onHandleConn4", "err", err.Error())
				}
				continue
			}
		}
	}
}

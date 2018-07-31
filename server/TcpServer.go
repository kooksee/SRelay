package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kooksee/srelay/types"
	"github.com/libp2p/go-reuseport"
	"github.com/patrickmn/go-cache"
)

var clientsCache *cache.Cache

type TcpServer struct {
}

func (ks *TcpServer) Listen(port int64) error {

	l, err := reuseport.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	go ks.accept(l)
	return nil
}

func (ks *TcpServer) onHandle(conn net.Conn) {
	kb := types.NewKBuffer()
	for {
		buf := make([]byte, 1024*16)
		n, err := conn.Read(buf)
		if err != nil {
			logger.Error("conn read error", "err", err.Error())

			if err == io.EOF {
				break
			}

			time.Sleep(time.Second * 2)
			continue
		}

		messages := kb.Next(buf[:n])
		if messages == nil {
			continue
		}

		for _, m := range messages {
			if m == nil || len(m) == 0 {
				continue
			}

			// 获得address 然后绑定客户端
			client, err := types.DecodeKClient(m)
			if err != nil {
				// 回复给客户端数据,json解析失败
				if _, err := conn.Write(types.ErrJsonParse(err)); err != nil {
					logger.Error("tcp onHandle 1", "err", err.Error())
				}
				continue
			}

			// 检查用户的签名
			if addr, b := cfg.CheckAddress(client.Sign); !b {
				if _, err := conn.Write(types.ErrSignError(errors.New(fmt.Sprintf("sign error")))); err != nil {
					logger.Error("tcp onHandle 3", "err", err.Error())
				}
			} else {
				clientsCache.SetDefault(addr, conn)
			}
		}
	}
}

func (ks *TcpServer) accept(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Error("tcp conn error ", "err", err)
			time.Sleep(time.Second * 3)
			continue
		}

		go ks.onHandle(c)
	}
}

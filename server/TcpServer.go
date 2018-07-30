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
			if !cfg.CheckAddress(client.ID, client.Sign) {
				if _, err := conn.Write(types.ErrSignError(errors.New(fmt.Sprintf("%s sign error", client.ID)))); err != nil {
					logger.Error("tcp onHandle 3", "err", err.Error())
				}
				continue
			}

			// 缓存，如果客户端没有定时确认连接，那么连接就会过时
			if !cfg.IsWhitelist(client.ID) {
				if _, err := conn.Write(types.ErrNotWhitelist(errors.New(fmt.Sprintf("%s not in white list", client.ID)))); err != nil {
					logger.Error("tcp onHandle 2", "err", err.Error())
				}
				continue
			}

			clientsCache.SetDefault(client.ID, conn)
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

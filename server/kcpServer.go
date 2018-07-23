package server

import (
	"bufio"
	"sync"
	"time"

	"github.com/kooksee/srelay/types"
	"github.com/kooksee/srelay/utils"
	knet "github.com/kooksee/srelay/utils/net"
)

var ksInstance *KcpServer
var ksOnce sync.Once

type KcpServer struct {
	l       *knet.KcpListener
	udps    map[int]*knet.UdpListener
	clients map[int]knet.Conn
}

func GetKcpServer() *KcpServer {
	ksOnce.Do(func() {
		ksInstance = &KcpServer{
		}
	})
	return ksInstance
}

func (ks *KcpServer) Listen() error {
	block, err := knet.GetCrypt(cfg.Crypt, cfg.Key, cfg.Salt)
	if err != nil {
		return err
	}

	ks.l, err = knet.ListenKcp(cfg.Host, cfg.KcpPort, block)
	if err != nil {
		return err
	}
	go ks.accept()
	return nil
}

// OpenPort 监听一个信的端口
func (ks *KcpServer) OpenPort(data []byte) error {
	clientInfo := new(types.ClientPort)
	if err := json.Unmarshal(data, &clientInfo); err != nil {
		return err
	}

	switch clientInfo.Protocol {
	case "http", "https":
	case "udp":
		usm := &UdpServerManager{}
		if err := usm.CreateUdp(int(clientInfo.Port)); err != nil {
			return err
		}
	}

	// 检查端口是否存在

	return nil
}

func (ks *KcpServer) onHandle(conn knet.Conn) {
	for {
		conn.LocalAddr().String()

		read := bufio.NewReader(conn)
		message, err := read.ReadBytes(types.Delim)
		if err != nil {
			logger.Info("kcp conn read message error", "err", err)
			break
		}

		message = utils.BytesTrimSpace(message)
		t := message[0]
		switch t {
		// 根据类型执行任务

		case 0x1:
			//   客户端请求打开端口
			//	 保存客户端的端口映射以及远程地址
			//	 保存客户端的名称

		case 0x2:
			//   客户端发送相应数据
			//   分析得到客户端地址
			//   根据地址找到连接
			//   然后把数据发送到未断开的连接上
			//   如果是http请求的话，那么就通过ID的方式发送给请求方
			//	 如果是udp请求，那么整个过程都是异步的，不要去等待数据
			//	 如果是ws请求，那么会有一个链接ID
			//	 缓存ID以及对应的结果数据

		}
		logger.Debug("message data", "data", string(message))
	}
}

func (ks *KcpServer) accept() {
	for {
		c, err := ks.l.Accept()
		if err != nil {
			logger.Error("kcp conn error ", "err", err)
			time.Sleep(time.Second * 3)
			continue
		}

		go ks.onHandle(c)
	}
}

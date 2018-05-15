package sp2p

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	kts "github.com/kooksee/srelay/types"
	knet "github.com/kooksee/srelay/utils/net"
)

func NewSP2p() *SP2p {
	p2p := &SP2p{
		nodeBackupT: time.NewTimer(10 * time.Minute),
		txC:         make(chan *kts.KMsg, 2000),
	}
	tab := newTable(PubkeyID(&p2p.priV.PublicKey), p2p.localAddr)
	p2p.tab = tab

	return p2p
}

type SP2p struct {
	tab *Table

	// 节点列表备份
	nodeBackupT *time.Timer
	priV        *ecdsa.PrivateKey
	txC         chan *kts.KMsg
	conn        knet.Conn
	localAddr   *net.UDPAddr
}

func (s *SP2p) Start() {
	block, err := knet.GetCrypt(cfg.Crypt, cfg.Key, cfg.Salt)
	if err != nil {
		panic(err.Error())
	}

	if c, err := knet.ConnectServer("kcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.KcpPort), block); err != nil {
		panic(err.Error())
	} else {
		s.conn = c
		go s.accept()
		s.handlerUdpPort()
		go s.handler()
		go s.ping()
	}
}

func (s *SP2p) Write(msg *kts.KMsg) error {
	return s.write(msg)
}

func (s *SP2p) WritePartition(msg *kts.KMsg) error {
	return s.write(msg)
}

func (s *SP2p) WriteBroadcast(msg *kts.KMsg) error {
	return s.write(msg)
}

func (s *SP2p) write(msg *kts.KMsg) error {
	if msg.FAddr == "" {
		msg.FAddr = s.localAddr.String()
	}
	if msg.FID == "" {
		msg.FID = s.tab.selfNode.ID.String()
	}
	if msg.ID == "" {
		msg.ID = string(UUID())
	}
	if msg.Version == "" {
		msg.Version = cfg.Version
	}
	if msg.TAddr == "" {
		return errors.New("目标地址不存在")
	}
	if _, err := s.conn.Write(msg.Dumps()); err != nil {
		return err
	}
	return nil
}

func (s *SP2p) handlerUdpPort() {
	t := time.NewTimer(2 * time.Second)
	msg := &kts.KMsg{
		Event: kts.CREATEUDPREQ,
		FAddr: s.localAddr.String(),
	}

	if err := s.write(msg); err != nil {
		logger.Error(err.Error())
	}

	for {
		select {
		case tx := <-s.txC:
			if tx.Event == kts.CREATEUDPRESP {
				if d, ok := tx.Data.(kts.Resp); ok {
					if d.Code == "ok" {
						s.localAddr = d.Data.(*net.UDPAddr)
						return
					} else {
						logger.Error(d.Msg)
					}
				}
			}
		case <-t.C:
			if err := s.write(msg); err != nil {
				logger.Error(err.Error())
				time.Sleep(time.Second)
				continue
			}
		}
	}
}

func (s *SP2p) handler() {
	for {
		tx := <-s.txC
		if hm.Contain(tx.Event) {
			go hm.GetHandler(tx.Event)(s, tx)
		}
	}
}

func (s *SP2p) accept() {
	s.conn.SetReadDeadline(time.Now().Add(connReadTimeout))
	read := bufio.NewReader(s.conn)
	for {
		message, err := read.ReadBytes(kts.Delim)
		if err != nil {
			logger.Info("kcp error ", err.Error())
			break
		}
		message = bytes.Trim(message, string(kts.Delim))
		logger.Debug("message", "data", string(message))

		msg := &kts.KMsg{}
		if err := msg.Decode(message); err != nil {
			logger.Error(err.Error())
			continue
		}
		s.txC <- msg
	}
}

func (s *SP2p) pingNode(taddr string) error {
	msg := &kts.KMsg{
		Event: "ping",
		TAddr: taddr,
	}
	return s.write(msg)
}

func (s *SP2p) findNode(taddr string, n int) error {
	msg := &kts.KMsg{
		Event: "findNode",
		TAddr: taddr,
		Data:  kts.FindNodeReq{NID: s.tab.selfNode.ID.String(), N: n},
	}
	return s.write(msg)
}

func (s *SP2p) ping() {
	tick := time.NewTicker(time.Minute)
	pingMsg := (&kts.KMsg{
		Event: kts.PINGREQ,
		FAddr: s.tab.GetNode().addr().String(),
	}).Dumps()

	for {
		select {
		case <-tick.C:
			s.conn.Write(pingMsg)
		}
	}
}

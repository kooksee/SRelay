package sp2p

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/dgraph-io/badger"
	kts "github.com/kooksee/srelay/types"
	knet "github.com/kooksee/srelay/utils/net"
)

func NewSP2p() *SP2p {
	p2p := &SP2p{
		nodeBackupTick: time.NewTicker(10 * time.Minute),
		pingTick:       time.NewTicker(10 * time.Minute),
		findNodeTick:   time.NewTicker(1 * time.Hour),
		txC:            make(chan *kts.KMsg, 2000),
	}
	tab := newTable(PubkeyID(&p2p.priV.PublicKey), p2p.localAddr)
	p2p.tab = tab

	go p2p.tickHandler()
	return p2p
}

type SP2p struct {
	tab *Table

	// 节点备份
	nodeBackupTick *time.Ticker
	// 节点ping
	pingTick *time.Ticker
	// 节点查询
	findNodeTick *time.Ticker

	priV      *ecdsa.PrivateKey
	txC       chan *kts.KMsg
	conn      knet.Conn
	localAddr *net.UDPAddr
}

func (s *SP2p) LoadSeeds(seeds []string) error {
	return cfg.Db.Update(func(txn *badger.Txn) error {
		k := []byte(cfg.NodesBackupKey)
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		for iter.Seek(k); ; iter.Next() {
			if !iter.ValidForPrefix(k) {
				break
			}

			val, err := iter.Item().Value()
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			seeds = append(seeds, string(val))
		}

		for _, rn := range seeds {
			n := MustParseNode(rn)
			n.updateAt = time.Now()
			s.tab.AddNode(n)
			go s.pingNode(n.addr().String())
		}

		// 节点启动的时候如果发现节点数量少,就去请求其他节点
		if s.tab.Size() < 3000 {
			// 每一个域选取一个节点
			for _, b := range s.tab.buckets {
				go s.findNode(b.Random().addr().String(), 8)
			}
		}

		return nil
	})
}

func (s *SP2p) dumpSeeds() {
	err := cfg.Db.Update(func(txn *badger.Txn) error {
		for _, n := range s.tab.GetAllNodes() {
			k := append([]byte(cfg.NodesBackupKey), n.ID.Bytes()...)
			if err := txn.Set(k, []byte(n.String())); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Error(err.Error())
	}
}

func (s *SP2p) tickHandler() {
	for {
		select {
		case <-s.nodeBackupTick.C:
			go s.dumpSeeds()
		case <-s.findNodeTick.C:
			for _, b := range s.tab.buckets {
				go s.findNode(b.Random().addr().String(), 8)
			}
		case <-s.pingTick.C:
			for _, n := range s.tab.FindRandomNodes(20) {
				go s.pingNode(n.addr().String())
			}
		}
	}
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
	s.conn.SetReadDeadline(time.Now().Add(cfg.ConnReadTimeout))
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

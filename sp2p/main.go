package sp2p

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/kooksee/srelay/protocol"
	kts "github.com/kooksee/srelay/types"
	knet "github.com/kooksee/srelay/utils/net"
)

func NewSP2p() *SP2p {
	p2p := &SP2p{
		txC: make(chan *kts.KMsg, 2000),
	}
	tab := newTable(PubkeyID(&cfg.PriV.PublicKey), p2p.localAddr)
	p2p.tab = tab

	p2p.localAddr = &net.UDPAddr{Port: cfg.Port, IP: net.ParseIP(cfg.Host)}
	return p2p
}

type SP2p struct {
	IP2p
	tab       *Table
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
		if s.tab.Size() < cfg.MinNodeSize {
			// 每一个域选取一个节点
			for _, b := range s.tab.buckets {
				b.peers.Each(func(index int, value interface{}) {
					go s.findNode(value.(*Node).addr().String(), 8)
				})
			}
		} else if s.tab.Size() < cfg.MaxNodeSize {
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
	pingMsg := (&kts.KMsg{
		Event: kts.PINGREQ,
		FAddr: s.tab.GetNode().addr().String(),
	}).Dumps()

	for {
		select {
		case <-cfg.NodeBackupTick.C:
			go s.dumpSeeds()
		case <-cfg.FindNodeTick.C:
			for _, b := range s.tab.buckets {
				go s.findNode(b.Random().addr().String(), 8)
			}
		case <-cfg.PingTick.C:
			for _, n := range s.tab.FindRandomNodes(20) {
				go s.pingNode(n.addr().String())
			}
		case <-cfg.PingKcpTick.C:
			s.conn.Write(pingMsg)
		}
	}
}

func (s *SP2p) StartP2p() {
	block, err := knet.GetCrypt(cfg.Crypt, cfg.Key, cfg.Salt)
	if err != nil {
		panic(err.Error())
	}

	if c, err := knet.ConnectServer("kcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), block); err != nil {
		panic(err.Error())
	} else {
		s.conn = c
		go s.accept()
		go s.handlerTx()
		go s.tickHandler()
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

func (s *SP2p) handlerTx() {
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

func (s *SP2p) GetTable() *Table {
	return s.tab
}
func (s *SP2p) pingNode(taddr string) error {
	return s.write(&kts.KMsg{
		Event: "ping",
		TAddr: taddr,
	})
}
func (s *SP2p) findNode(taddr string, n int) error {
	return s.write(&kts.KMsg{
		Event: "findNode",
		TAddr: taddr,
		Data:  protocol.FindNodeReq{NID: s.tab.selfNode.ID.String(), N: n},
	})
}

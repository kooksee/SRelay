package sp2p

import (
	"crypto/ecdsa"
	"net"
	"time"

	"github.com/kooksee/uspnet/log"

	knet "github.com/kooksee/srelay/utils/net"
	"github.com/kooksee/uspnet/net/netutil"
	"github.com/vmihailenco/msgpack"
)

const (
	nodesBackupTime = 10 * time.Minute
)

// udp implements the RPC protocol.
type KRpc struct {
	IKRpc

	netRestrict *netutil.Netlist
	localAddr   *net.UDPAddr

	sendChan chan *sendQ
	recvChan chan *recvQ

	tab  *Table
	conn knet.Conn
}

// pending represents a pending reply.
//
// some implementations of the protocol wish to send more than one
// reply packet to findnode. in general, any neighbors packet cannot
// be matched up with a specific findnode packet.
//
// our implementation handles this by storing a callback function for
// each pending reply. incoming packets from a node are dispatched
// to all the callback functions for that node.

type sendQ struct {
	tx       *Tx
	deadline time.Time
}

type recvQ struct {
	tx       *Tx
	deadline time.Time
	buf      []byte
}

// ListenUDP returns a new table that listens for UDP packets on laddr.
func NewKRpc(priV *ecdsa.PrivateKey, conn knet.Conn, nodeList []string, netRestrict *netutil.Netlist) (*KRpc, error) {

	kRpc := &KRpc{
		netRestrict: netRestrict,
		sendChan:    make(chan *sendQ, 10000),
		recvChan:    make(chan *recvQ, 10000),
		conn:        conn,
	}

	kRpc.localAddr = kRpc.conn.LocalAddr().(*net.UDPAddr)
	tab := newTable(PubkeyID(&priV.PublicKey), kRpc.localAddr)

	// 加载节点列表
	for _, n := range nodeList {
		tab.AddNode(MustParseNode(n))
	}
	kRpc.tab = tab

	go kRpc.loop()
	go kRpc.readLoop()
	go kRpc.writeLoop()

	log.Info("UDP listener up", "self", kRpc.tab.selfNode)
	return kRpc, nil
}


// loop runs in its own goroutine. it keeps track of
// the refresh timer and the pending reply queue.
func (t *KRpc) loop() {
	var (
		timeout      = time.NewTimer(0)
		contTimeouts = 0 // number of continuous timeouts to do NTP checks
		ntpWarnTime  = time.Unix(0, 0)
	)

	handleTx := func(buf []byte, tx *Tx) error {
		if err := msgpack.Unmarshal(buf, tx); err != nil {
			return err
		}

		return tx.Packet.OnHandle(t, tx)
	}

	<-timeout.C // ignore first timeout
	defer timeout.Stop()

	for {
		select {
		case p := <-t.recvChan:
			if expired(p.deadline.Unix()) {
				continue
			}

			if err := handleTx(p.buf, p.tx); err != nil {
				log.Error(err.Error())
				t.recvChan <- p
			}

		case <-timeout.C:
			contTimeouts++

			if contTimeouts > ntpFailureThreshold {
				if time.Since(ntpWarnTime) >= ntpWarningCooldown {
					ntpWarnTime = time.Now()
					go checkClockDrift()
				}
				contTimeouts = 0
			}
		}
	}
}

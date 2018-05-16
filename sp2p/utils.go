package sp2p

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/kooksee/srelay/types"
	"github.com/kooksee/uspnet/common"
	"github.com/kooksee/uspnet/crypto/secp256k1"
	"github.com/satori/go.uuid"
)

// recoverNodeID computes the public key used to sign the
// given hash from the signature.
func RecoverNodeID(hash, sig []byte) (id NodeID, err error) {
	pubKey, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		return id, err
	}
	if len(pubKey)-1 != len(id) {
		return id, fmt.Errorf("recovered pubkey has %d bits, want %d bits", len(pubKey)*8, (len(id)+1)*8)
	}
	for i := range id {
		id[i] = pubKey[i+1]
	}
	return id, nil
}

// DistCmp compares the distances a->target and b->target.
// Returns -1 if a is closer to target, 1 if b is closer to target
// and 0 if they are equal.
func distCmp(target, a, b common.Hash) int {
	for i := range target {
		da := a[i] ^ target[i]
		db := b[i] ^ target[i]
		if da > db {
			return 1
		} else if da < db {
			return -1
		}
	}
	return 0
}

func UUID() []byte {
	for {
		uid, err := uuid.NewV4()
		if err == nil {
			return uid.Bytes()
		}
	}
	return nil
}

func expired(ts int64) bool {
	return time.Unix(ts, 0).Before(time.Now())
}

func timeAdd(ts time.Duration) time.Time {
	return time.Now().Add(ts)
}

type KBuffer struct {
	buf   map[string][]byte
	Delim []byte

	sync.RWMutex
}

func (t *KBuffer) Add(key string, b []byte) {
	t.Lock()
	defer t.Unlock()

	t.buf[key] = append(t.buf[key], b...)
}

func (t *KBuffer) Next(key string) [][]byte {
	t.RLock()
	defer t.RUnlock()

	if len(t.buf) != 0 {
		d := bytes.Split(t.buf[key], t.Delim)
		if len(d) > 1 {
			t.buf[key] = d[len(d)-1]
			return d[:len(d)-2]
		}
	}
	return nil
}

func If(cond bool, trueVal, falseVal interface{}) interface{} {
	if cond {
		return trueVal
	}
	return falseVal
}

// logdist returns the logarithmic distance between a and b, log2(a ^ b).
func logdist(a, b common.Hash) int {
	lz := 0
	for i := range a {
		x := a[i] ^ b[i]
		if x == 0 {
			lz += 8
		} else {
			lz += lzcount[x]
			break
		}
	}
	return len(a)*8 - lz
}

// hashAtDistance returns a random hash such that logdist(a, b) == n
func hashAtDistance(a common.Hash, n int) (b common.Hash) {
	if n == 0 {
		return a
	}
	// flip bit at position n, fill the rest with random bits
	b = a
	pos := len(a) - n/8 - 1
	bit := byte(0x01) << (byte(n%8) - 1)
	if bit == 0 {
		pos++
		bit = 0x80
	}
	b[pos] = a[pos]&^bit | ^a[pos]&bit // TODO: randomize end bits
	for i := pos + 1; i < len(a); i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}

func randUint(max uint32) uint32 {
	if max == 0 {
		return 0
	}
	rand.Seed(time.Now().Unix())
	return rand.Uint32() % max
}

func NodeFromKMsg(msg *types.KMsg) *Node {
	nid, err := HexID(msg.FID)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	addr, err := net.ResolveUDPAddr("udp", msg.FAddr)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return NewNode(nid, addr.IP, uint16(addr.Port))
}
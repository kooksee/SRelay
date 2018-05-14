package tclient

import (
	"net/http"
	"time"

	klog "github.com/sirupsen/logrus"

	knet "SRelay/utils/net"

	"github.com/gorilla/websocket"
)

var (
	tcpClients map[string]knet.Conn
	wsClients  map[string]*websocket.Conn
	log        *klog.Entry
	upgrader   = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	connReadTimeout time.Duration = 10 * time.Second
)
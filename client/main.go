package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kooksee/srelay/types"
	"github.com/satori/go.uuid"
)

type config struct {
	addr string
}

func main() {

	fmt.Println("hello")

	servertAddr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
	if err != nil {
		panic(err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, servertAddr)
	if err != nil {
		panic(err.Error())
	}

	Id := uuid.Must(uuid.NewV4()).String()
	conn.Write(types.KMsg{ID: Id}.Bytes())

	go func() {

		for {
			buf := make([]byte, 1024*16)
			_, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}

				time.Sleep(time.Second * 2)
				continue
			}

			fmt.Println(string(types.BytesTrimSpace(buf)))
		}

		time.Sleep(time.Second * 2)

	}()

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		panic(err.Error())
	}

	uconn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err.Error())
	}

	go func() {

		for {
			uconn.Write(types.KMsg{
				ID:   Id,
				Data: []byte("ok"),
			}.Bytes())

			//time.Sleep(time.Microsecond)
		}

	}()

	select {}

}

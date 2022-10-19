package network

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

var MyAddress string
var readWriteDdl time.Duration = 5 * 10^9

func listenTcp() error {
	listener, err := net.Listen("tcp", MyAddress)
	if err != nil {
		log.Printf(err.Error())
	}
	defer func() { _ = listener.Close() }()
	log.Printf("bound to &q", listener.Addr())
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf(err.Error())
			return err
		}
		err = conn.SetDeadline(time.Now().Add(readWriteDdl))
		if err != nil {
			log.Printf(err.Error())
		}
		go func(c net.Conn) {
			defer func() {
				c.Close()
			}()

			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Println(err.Error())
				return
			}
			temp := strings.TrimSpace(string(netData))
			if temp == "STOP" {
				return
			}
			log.Printf("received: %q", temp)
		}(conn)
	}
}

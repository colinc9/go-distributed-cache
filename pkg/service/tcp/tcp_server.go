package tcp

import (
	"github.com/colinc9/go-distributed-cache/pkg/service"
	"log"
	"net"
	"time"
)

var MyAddress string
var readWriteDdl time.Duration = 5 * time.Minute

func ListenTcp() error {
	listener, err := net.Listen("tcp", MyAddress)
	if err != nil {
		log.Printf(err.Error())
	}
	defer func() { _ = listener.Close() }()
	log.Printf("bound to %+v", listener.Addr())
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
			msg, err := Decoder(c)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Printf("received: %+v", *msg)
			service.HandelMsg(msg)
		}(conn)
	}
}

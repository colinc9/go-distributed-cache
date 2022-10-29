package network

import (
	"log"
	"net"
	"time"
)

var TargetAddress []string
var timeLimit time.Duration = 5 * time.Minute // timeout limit in nano sec

func DialTcp() {
	for _, target := range TargetAddress {
		go func(target string) {
			conn, err := net.DialTimeout("tcp", target, timeLimit)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			defer conn.Close()

			for {
				message := "sent message to " + target
				log.Print("->: " + message)
				conn.Write([]byte(message))

				time.Sleep(5 * time.Second)
			}
		}(target)

	}

}
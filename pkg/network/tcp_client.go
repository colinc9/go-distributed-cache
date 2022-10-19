package network

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var TargetAddress string
var timeLimit time.Duration = 5 * 10^9 // timeout limit in nano sec

func dialTcp() {
	conn, err := net.DialTimeout("tcp", TargetAddress, timeLimit)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		log.Print(">> ")
		text, _ := reader.ReadString('\n')
		log.Printf(text+"\n")

		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("->: " + message)
		conn.Write([]byte(message))

		if strings.TrimSpace(string(text)) == "STOP" {
			log.Println("TCP client exiting...")
			return
		}
	}
}
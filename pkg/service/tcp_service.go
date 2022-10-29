package service

import (
	"errors"
	"github.com/colinc9/go-distributed-cache/pkg/model"
	"github.com/colinc9/go-distributed-cache/pkg/service/tcp"
	"log"
)

type TCPService struct {
	Cache *model.LRUCache
}

var TcpService *TCPService

func HandelMsg(msg *tcp.Message) error {
	var err error = nil
	switch msg.Type {
		case tcp.Get:
			_, ok := TcpService.Cache.Get(msg.Key)
			if !ok {
				err = errors.New("Get Failed, Key: " + msg.Key.(string))
				log.Printf(err.Error())
			}
		case tcp.Set:
			_, _, ok  := TcpService.Cache.Set(msg.Key, msg.Value)
			if !ok {
				err = errors.New("Set Failed, Key Value: " + msg.Key.(string) + " " + msg.Value.(string))
				log.Printf(err.Error())
			}
	}
	return err
}

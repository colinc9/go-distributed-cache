package controller

import (
	"github.com/colinc9/go-distributed-cache/pkg/config"
	"github.com/colinc9/go-distributed-cache/pkg/model"
	"github.com/colinc9/go-distributed-cache/pkg/service"
	"github.com/colinc9/go-distributed-cache/pkg/service/tcp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var cacheService *service.CacheService

func Run() error {
	server := setGinRouter()
	go func() {
		tcp.ListenTcp()
	}()
	LRUcache, error := model.NewLRUCache(10)
	if error != nil{
		cacheService = &service.CacheService{
			Cache: LRUcache,
		}
		tcp.TcpService = &tcp.TCPService{
			Cache: LRUcache,
		}
	}
	server.Run(config.GetDefaultInsCfg().AppAddress)

	return nil
}

func setGinRouter() *gin.Engine {
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()
	router.GET("/", HealthCheck)
	router.GET("/get/:key", Get)
	router.POST("/set/:key/value/:value",Set)
	router.GET("/test", SendMsgToTask)
	return router
}

func Get(c *gin.Context) {
	value, ok := cacheService.Get(c.Param("key"))
	if ok {
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusBadRequest, c.Param("key"))
	}

}

func Set(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")

	_, _, ok := cacheService.Set(key, value)

	if ok {
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusBadRequest, c.Param("key"))
	}
}

func HealthCheck(c *gin.Context) {
	log.Printf("server is listening...")
	c.IndentedJSON(http.StatusOK, "Alive!")
}

func SendMsgToTask(c *gin.Context) {
	msg := tcp.Message{Type: tcp.Test}
	tcp.DialTcp(&msg)
	c.IndentedJSON(http.StatusOK, "Alive!")
}







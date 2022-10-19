package controller

import (
	"github.com/colinc9/go-distributed-cache/pkg/config"
	"github.com/colinc9/go-distributed-cache/pkg/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var cache *model.LRUCache

func Run() error {
	server := setGinRouter()
	server.Run(config.GetDefaultInsCfg().AppAddress)

	return nil
}

func setGinRouter() *gin.Engine {
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()
	LRUcache, error := model.NewLRUCache(10)
	if error != nil{
		cache = LRUcache
	}
	router.GET("/", HealthCheck)
	router.GET("/get/:key", Get)
	router.POST("/set/:key/value/:value",Set)
	return router
}

func Get(c *gin.Context) {
	value, ok := cache.Get(c.Param("key"))
	if ok {
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusBadRequest, c.Param("key"))
	}

}

func Set(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")

	_, _, ok := cache.Set(key, value)

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







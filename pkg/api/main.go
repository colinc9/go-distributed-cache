package api

import (
	"github.com/colinc9/go-distributed-cache/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var cache *model.LRUCache

func main() {
	router := gin.Default()
	LRUcache, error := model.NewLRUCache(10)
	if error != nil{
		cache = LRUcache
	}
	router.GET("/", healthCheck)
	router.GET("/get/:key", get)
	router.POST("/set/:key/value/:value", set)
	router.Run("0.0.0.0:8080")
}

func get(c *gin.Context) {
	value, ok := cache.Get(c.Param("key"))
	if ok {
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusBadRequest, c.Param("key"))
	}

}

func set(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")

	_, _, ok := cache.Set(key, value)

	if ok {
		c.IndentedJSON(http.StatusOK, value)
	} else {
		c.IndentedJSON(http.StatusBadRequest, c.Param("key"))
	}
}

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Alive!")
}


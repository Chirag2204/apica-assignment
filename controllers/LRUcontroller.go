package controllers

import (
	"lru-cache-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var cache = models.NewLRUCache(5)

type CacheRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func GetCache(c *gin.Context) {
	key := c.Param("key")
	value, exists := cache.Get(key)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func SetCache(c *gin.Context) {
	key := c.Param("key")
	var request CacheRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cache.Set(key, request.Value)
	c.JSON(http.StatusOK, gin.H{"message": "Key set successfully"})
}

func DeleteCache(c *gin.Context) {
	key := c.Param("key")
	cache.Delete(key)
	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}

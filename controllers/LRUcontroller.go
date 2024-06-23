package controllers

import (
	"lru-cache-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var cache = models.NewLRUCache(5)

func GetCache(c *gin.Context) {
	key := c.Param("key")
	if value, ok := cache.Get(key); ok {
		c.JSON(http.StatusOK, gin.H{"value": value})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
	}
}

func SetCache(c *gin.Context) {
	key := c.Param("key")
	var json struct {
		Value string `json:"value"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cache.Set(key, json.Value)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DeleteCache(c *gin.Context) {
	key := c.Param("key")
	cache.Delete(key)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

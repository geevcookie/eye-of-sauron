package handlers

import (
	"eye-of-sauron/internal/metrics/collector"
	"github.com/gin-gonic/gin"
)

func Ping(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
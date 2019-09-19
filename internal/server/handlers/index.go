package handlers

import (
	"net/http"

	"eye-of-sauron/internal/metrics/collector"
	"github.com/gin-gonic/gin"
)

func Index(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{})
	}
}
package handlers

import (
	"encoding/json"
	"net/http"

	"eye-of-sauron/internal/metrics/collector"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// GetMetrics godoc
// @Summary Get all metrics
// @Description Returns a JSON object containing all metrics for CPU, Memory, and Disk
// @Accept json
// @Produce json
// @Success 200 {object} collector.Metrics "JSON object containing all the metrics"
// @tags metrics
// @Router /api/metrics [get]
func GetMetrics(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, collector.Stats.Metrics())
	}
}

// GetCPU godoc
// @Summary Get CPU metrics
// @Description Returns a JSON object containing all CPU metrics
// @Accept json
// @Produce json
// @Success 200 {object} collector.CPUMetrics "JSON object containing CPU metrics"
// @tags metrics
// @Router /api/metrics/cpu [get]
func GetCPU(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, collector.Stats.Metrics().CPU)
	}
}

// GetMemory godoc
// @Summary Get CPU metrics
// @Description Returns a JSON object containing all memory metrics
// @Accept json
// @Produce json
// @Success 200 {object} collector.MemoryMetrics "JSON object containing memory metrics"
// @tags metrics
// @Router /api/metrics/memory [get]
func GetMemory(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, collector.Stats.Metrics().Memory)
	}
}

// GetLoad godoc
// @Summary Get load metrics
// @Description Returns a JSON object containing all load metrics
// @Accept json
// @Produce json
// @Success 200 {object} collector.LoadMetrics "JSON object containing load metrics"
// @tags metrics
// @Router /api/metrics/load [get]
func GetLoad(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, collector.Stats.Metrics().Load)
	}
}

// GetDisks godoc
// @Summary Get disk metrics
// @Description Returns a JSON object containing all disk metrics
// @Accept json
// @Produce json
// @Success 200 {array} collector.DiskMetrics "JSON array containing disk metrics"
// @tags metrics
// @Router /api/metrics/disks [get]
func GetDisks(collector *collector.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, collector.Stats.Metrics().Disks)
	}
}

func WrapWebSocket(channel *chan collector.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		wshandler(c.Writer, c.Request, channel)
	}
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request, c *chan collector.Metrics) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("Failed to set websocket upgrade!")
		return
	}

	for {
		metrics := <-*c

		msg, _ := json.Marshal(metrics)
		_ = conn.WriteMessage(1, msg)
	}
}

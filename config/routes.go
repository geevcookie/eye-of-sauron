package config

import (
	"eye-of-sauron/internal/metrics/collector"
	"eye-of-sauron/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

// Verb represents a HTTP Verb
type Verb int

// Define all available verbs
const (
	GET    Verb = 0
	POST   Verb = 1
	PUT    Verb = 2
	PATCH  Verb = 3
	DELETE Verb = 4
)

type HandlerFunc func(*collector.Collector) gin.HandlerFunc

// Route represents a single route in the system
type Route struct {
	Verb       Verb
	Path       string
	Handler    HandlerFunc
}

// Routes returns a []Route representing all the available routes in the API
func Routes() []Route {
	// Define routes here
	return []Route{
		{GET, "/ping", handlers.Ping},
		{GET, "/", handlers.Index},
		{GET, "/api/metrics", handlers.GetMetrics},
		{GET, "/api/metrics/cpu", handlers.GetCPU},
		{GET, "/api/metrics/memory", handlers.GetMemory},
		{GET, "/api/metrics/disks", handlers.GetDisks},
		{GET, "/api/metrics/load", handlers.GetLoad},
	}
}

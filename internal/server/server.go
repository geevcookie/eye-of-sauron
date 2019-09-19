package server

import (
	"fmt"
	"time"

	"eye-of-sauron/config"
	"eye-of-sauron/internal"
	"eye-of-sauron/internal/metrics/collector"
	"eye-of-sauron/internal/server/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"

	_ "eye-of-sauron/docs"
)

// Server represents the API server
type Server struct {
	Router *gin.Engine
}

// NewServer creates a new Server
func NewServer() Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(ginlogrus.Logger(log.New()))
	r.Use(gin.Recovery())

	// Asset and HTML Config
	r.Static("/assets", "./static")
	r.LoadHTMLGlob("internal/templates/*.gohtml")

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// Swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return Server{Router: r}
}

// Run starts the server on the configured port
func (s *Server) Run(cfg internal.Config) error {
	log.Infof("Starting server on 127.0.0.1:%s...", cfg.ServerConfig.Port)
	return s.Router.Run(fmt.Sprintf(":%s", cfg.ServerConfig.Port))
}

// LoadRoutes adds all the required routes to the API server
func (s *Server) LoadRoutes(routes []config.Route, c *collector.Collector, channel *chan collector.Metrics) {
	// WebSocket route
	s.Router.GET("/ws", handlers.WrapWebSocket(channel))

	// Process routes from configuration
	for _, r := range routes {
		switch r.Verb {
		case config.GET:
			s.Router.GET(r.Path, r.Handler(c))
		case config.POST:
			s.Router.POST(r.Path, r.Handler(c))
		case config.PUT:
			s.Router.PUT(r.Path, r.Handler(c))
		case config.PATCH:
			s.Router.PATCH(r.Path, r.Handler(c))
		case config.DELETE:
			s.Router.DELETE(r.Path, r.Handler(c))
		}
	}
}

package router

import (
	"fmt"
	"time"

	"notes-app/server/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notes_api_requests_total",
			Help: "Total HTTP requests by method, path and status",
		},
		[]string{"method", "path", "status"},
	)
	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "notes_api_request_duration_seconds",
			Help:    "Request duration in seconds by method and path",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func Setup(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(reqDuration)

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		reqCount.WithLabelValues(c.Request.Method, path, fmt.Sprintf("%d", c.Writer.Status())).Inc()
		reqDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	})

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	noteHandler := handlers.NewNoteHandler(db)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/api/health", handlers.HealthCheck)

	api := r.Group("/api/notes")
	{
		api.GET("", noteHandler.List)
		api.GET("/:id", noteHandler.Get)
		api.POST("", noteHandler.Create)
		api.PUT("/:id", noteHandler.Update)
		api.DELETE("/:id", noteHandler.Delete)
	}

	return r
}

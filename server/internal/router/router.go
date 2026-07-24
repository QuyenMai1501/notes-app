package router

import (
	"notes-app/server/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

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

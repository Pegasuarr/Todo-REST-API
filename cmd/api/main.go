package main

import (
	"Todo_App/internal/config"
	"Todo_App/internal/database"
	"Todo_App/internal/handlers"
	"Todo_App/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.Db)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()

	// Enable CORS for API flexibility
	router.Use(middleware.CORSMiddleware())

	// Root healthcheck endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo API is running!!!!",
			"status":   "success",
			"database": "connected",
		})
	})

	// Todo REST API endpoints
	router.GET("/todos", handlers.GetTodosHandler(pool))
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))
	router.PUT("/todos/:id", handlers.UpdateTodoHandler(pool))
	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))

	log.Println("Starting server on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

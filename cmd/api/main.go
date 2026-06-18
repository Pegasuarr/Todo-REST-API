package main

import (
	"Todo_App/internal/config"
	"Todo_App/internal/database"
	"log"

	"Todo_App/internal/handlers"
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
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			//map [string] of interface
			//map [string] any {}
			"message":  "Todo API is running!!!!",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler(pool))

	if err := router.Run(":" + cfg.Port); err != nil {
	log.Fatal("Failed to start server:", err)
}

}

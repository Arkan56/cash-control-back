package main

import (
	"cash-control/internal/config"
	"cash-control/internal/database"
	"cash-control/internal/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failde to connect to database:", err)
	}

	defer pool.Close()

	fmt.Println("API Server")
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":  "it's ok",
			"databese": "connected",
		})
	})

	router.POST("/stores", handlers.CreateStoreHandler(pool))
	router.GET("/stores", handlers.GetAllStoresHandler(pool))

	router.Run(":" + cfg.Port)

}

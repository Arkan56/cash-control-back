package main

import (
	"cash-control/internal/config"
	"cash-control/internal/database"
	"cash-control/internal/handlers"
	"cash-control/internal/middleware"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	const (
		admin_rol  = 1
		worker_rol = 2
	)
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":  "it's ok",
			"databese": "connected",
		})
	})

	router.POST("/auth/login/user", handlers.LoginUserHandler(pool, cfg))

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))

	admin := api.Group("/admin")
	admin.Use(middleware.AuthRolMiddleware(admin_rol))
	{
		admin.POST("/stores", handlers.CreateStoreHandler(pool))
		admin.POST("/vaults", handlers.CreateVaultHandler(pool))
		admin.POST("/users", handlers.CreateUserHandler(pool))
	}

	core := api.Group("/core")
	core.Use(middleware.AuthRolMiddleware(admin_rol, worker_rol))
	{
		core.GET("/stores", handlers.GetAllStoresHandler(pool))
		core.POST("/movements", handlers.CreateMovHandler(pool))
		core.GET("/movements/vault/:id", handlers.GetAllMovsHandler(pool))
		core.GET("/vaults/:id", handlers.GetAllVaultsByStoreIdHandler(pool))
		core.POST("/auth/login/vault", handlers.LoginVaultHandler(pool))
	}
	router.Run(":" + cfg.Port)

}

package handlers

import (
	"cash-control/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateStoreInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateStoreHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateStoreInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		store, err := repository.CreateStore(pool, input.Name)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, store)
	}
}

func GetAllStoresHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		stores, err := repository.GetAllStores(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stores)
	}
}

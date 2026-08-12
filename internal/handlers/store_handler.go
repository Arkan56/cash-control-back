package handlers

import (
	"cash-control/internal/repository"
	"net/http"
	"strconv"

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

func GetStoresByUserId(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			return
		}

		userIdFl := userIdInterface.(float64)

		userId := int64(userIdFl)

		stores, err := repository.GetStoresByUserId(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stores)
	}
}

func CreateStoreAccessHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			StoreID int64 `json:"store_id"`
			UserID  int64 `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := repository.CreateStoreAccess(pool, request.UserID, request.StoreID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "store access created successfully",
		})
	}
}

func GetStoresAccesByUserId(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id" + err.Error()})
			return
		}

		storesAccess, err := repository.GetStoresAccessByUserId(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, storesAccess)
	}
}

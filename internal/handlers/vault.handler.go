package handlers

import (
	"cash-control/internal/models"
	"cash-control/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateVaultHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateVaultRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vault := &models.CreateVaultRequest{
			Name:    input.Name,
			StoreId: input.StoreId,
		}

		createdVault, err := repository.CreateVault(pool, vault)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdVault)
	}
}

func GetAllVaultsByStoreIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeId, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store id" + err.Error()})
			return
		}

		vaults, err := repository.GetVaultsByStoreId(pool, storeId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vaults)
	}
}

func GetVaultsByUserId(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			return
		}

		userIdFl := userIdInterface.(float64)

		userId := int64(userIdFl)

		vaults, err := repository.GetVaultsByUserId(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vaults)
	}
}

func CreateVaultAccessHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			VaultID int64 `json:"vault_id"`
			UserID  int64 `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := repository.CreateVaultAccess(pool, request.UserID, request.VaultID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "vault access created successfully",
		})
	}
}

func GetVaultById(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		vaultId, err := strconv.ParseInt(c.Param("vaultId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault id" + err.Error()})
			return
		}

		vault, err := repository.GetVaultById(pool, vaultId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vault)
	}
}

func GetVaultsAccesByUserId(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id" + err.Error()})
			return
		}

		vaultsAccess, err := repository.GetVaultAccessByUserId(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vaultsAccess)
	}
}

func GetAllVaultsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		vaults, err := repository.GetAllVaults(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vaults)
	}
}

func GetVaultsByStoreUserId(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			return
		}

		userIdFl := userIdInterface.(float64)

		userId := int64(userIdFl)

		storeId, err := strconv.ParseInt(c.Param("storeId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault id" + err.Error()})
			return
		}

		vaults, err := repository.GetVaultsByStoreUserId(pool, userId, storeId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vaults)
	}
}

package handlers

import (
	"cash-control/internal/models"
	"cash-control/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type LoginVaultReq struct {
	Id       int64  `json:"vault_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginVaultResponse struct {
	ID        int64  `json:"id"`
	StoreName string `json:"store_name"`
	Name      string `json:"name"`
	Balance   int64  `json:"balance"`
}

func CreateVaultHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateVaultRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password" + err.Error()})
			return
		}

		vault := &models.CreateVaultRequest{
			Name:     input.Name,
			Password: string(hashedPassword),
			StoreId:  input.StoreId,
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

func LoginVaultHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginVaultReq

		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vault, err := repository.GetVaultById(pool, loginRequest.Id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(vault.Password), []byte(loginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		response := LoginVaultResponse{
			ID:        vault.ID,
			StoreName: vault.StoreName,
			Name:      vault.Name,
			Balance:   vault.Balance,
		}

		c.JSON(http.StatusOK, response)
	}
}

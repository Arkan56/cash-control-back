package handlers

import (
	"cash-control/internal/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateMovInput struct {
	Detail           string `json:"detail"  binding:"required"`
	Amount           int64  `json:"amount"  binding:"required"`
	AmountCategoryID int32  `json:"amount_category_id" binding:"required"`
	VaultID          int64  `json:"vault_id"  binding:"required"`
	UserID           int64  `json:"user_id"  binding:"required"`
}

func CreateMovHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateMovInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		movement, err := repository.CreateMovement(pool,
			input.Detail,
			input.Amount,
			input.AmountCategoryID,
			input.VaultID,
			input.UserID)

		if err != nil {
			fmt.Printf("Error detallado de DB: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, movement)
	}
}
func GetAllMovsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		vaultID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid vault id"})
			return
		}

		movs, err := repository.GetAllMovements(pool, vaultID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movs)
	}
}

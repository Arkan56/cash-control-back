package handlers

import (
	"cash-control/internal/repository"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateMovInput struct {
	Detail           string `json:"detail"  binding:"required"`
	Amount           int64  `json:"amount"  binding:"required"`
	AmountCategoryID int32  `json:"amount_category_id" binding:"required"`
	VaultID          int64  `json:"vault_id"  binding:"required"`
}

func CreateMovHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateMovInput
		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			return
		}

		userIdFl := userIdInterface.(float64)

		userId := int64(userIdFl)

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		movement, err := repository.CreateMovement(pool,
			input.Detail,
			input.Amount,
			input.AmountCategoryID,
			input.VaultID,
			userId)

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
		vaultID, err := strconv.ParseInt(c.Param("vaultId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault id"})
			return
		}

		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
			return
		}

		if endDate.Before(startDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end date cannot be before start date"})
			return
		}

		endDate = endDate.AddDate(0, 0, 1)

		movs, err := repository.GetAllMovements(
			pool,
			vaultID,
			startDate,
			endDate,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movs)
	}
}

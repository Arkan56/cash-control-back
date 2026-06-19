package handlers

import (
	"cash-control/internal/config"
	"cash-control/internal/models"
	"cash-control/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateUserRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password" + err.Error()})
			return
		}

		user := &models.CreateUserRequest{
			UserName: input.UserName,
			Name:     input.Name,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateWorkerUser(pool, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, createdUser)
	}
}

func LoginUserHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginUserRequest

		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByUserName(pool, loginRequest.UserName)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bad credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bad credentials"})
			return
		}

		claims := jwt.MapClaims{
			"user_id":     user.ID,
			"user_name":   user.UserName,
			"user_rol_id": user.IdRol,
			"exp":         time.Now().Add(1 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Feiled to generate the token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.LoginUserResponse{Token: tokenString})

	}
}

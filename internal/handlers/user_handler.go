package handlers

import (
	"cash-control/internal/config"
	"cash-control/internal/models"
	"cash-control/internal/repository"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SyncAccessRequest struct {
	StoreIDs []int64 `json:"store_ids"`
	VaultIDs []int64 `json:"vault_ids"`
}

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
			IdRol:    input.IdRol,
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

func GetAllUsersHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := repository.GetAllUsers(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUserByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id" + err.Error()})
			return
		}

		user, err := repository.GetUserById(pool, userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SyncUserAccessHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido"})
			return
		}

		var req SyncAccessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.StoreIDs == nil {
			req.StoreIDs = []int64{}
		}
		if req.VaultIDs == nil {
			req.VaultIDs = []int64{}
		}

		success, err := repository.SyncAccessByUserId(pool, userId, req.StoreIDs, req.VaultIDs)

		if err != nil || !success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "synchronized access"})
	}
}

func UpdateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var req models.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.UserName == nil && req.Name == nil && req.Password == nil && req.IdRol == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
			return
		}

		var hashed *string
		if req.Password != nil {
			h, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}
			s := string(h)
			hashed = &s
		}

		user, err := repository.UpdateUser(pool, userId,
			req.UserName, req.Name, hashed, req.IdRol)
		if err != nil {
			var pgErr *pgconn.PgError
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			case errors.As(err, &pgErr) && pgErr.Code == "23505":
				c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			case errors.As(err, &pgErr) && pgErr.Code == "23503":
				c.JSON(http.StatusBadRequest, gin.H{"error": "role does not exist"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

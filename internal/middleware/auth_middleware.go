package middleware

import (
	"cash-control/internal/config"
	"cash-control/internal/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Claims"})
			c.Abort()
			return
		}

		user_rol_id, ok := claims["user_rol_id"].(float64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Claims"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)

			if time.Now().After(expirationTime) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userID)
		c.Set("user_rol_id", user_rol_id)
		c.Next()
	}
}

func AuthRolMiddleware(roles ...int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIdInterface, exist := c.Get("user_rol_id")

		if !exist {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_role_id not found in context"})
			c.Abort()
			return
		}

		roleId := roleIdInterface.(float64)

		for _, allowedRole := range roles {
			if int32(roleId) == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}

func VaultAccessMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("VaultAccessMiddleware:", c.Request.URL.Path)
		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			c.Abort()
			return
		}

		userIdFl := userIdInterface.(float64)
		userId := int64(userIdFl)

		userRolIdInterface, exist := c.Get("user_rol_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_rol_id not found"})
			c.Abort()
			return
		}

		userRolIdFl := userRolIdInterface.(float64)
		userRolId := int32(userRolIdFl)

		if userRolId == 1 {
			c.Next()
			return
		}

		vaultId, err := strconv.ParseInt(c.Param("vaultId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault id" + err.Error()})
			c.Abort()
			return
		}

		var userAcces bool
		userAcces, err = repository.UserHasVaultAcces(pool, userId, vaultId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !userAcces {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acces denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func StoreAccessMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIdInterface, exist := c.Get("user_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
			return
		}

		userIdFl := userIdInterface.(float64)
		userId := int64(userIdFl)

		userRolIdInterface, exist := c.Get("user_rol_id")

		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_rol_id not found"})
			return
		}

		userRolIdFl := userRolIdInterface.(float64)
		userRolId := int32(userRolIdFl)

		if userRolId == 1 {
			c.Next()
			return
		}

		storeId, err := strconv.ParseInt(c.Param("storeId"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store id" + err.Error()})
			return
		}

		var userAcces bool
		userAcces, err = repository.UserHasStoreAcces(pool, userId, storeId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !userAcces {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acces denied"})
			return
		}

		c.Next()
	}
}

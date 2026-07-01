package handlers

import (
	"Todo_App/internal/auth"
	"Todo_App/internal/repository"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := repository.CreateUser(pool, input.Username, string(hashedBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Registration successful!",
			"user":    user,
		})
	}
}

func LoginHandler(pool *pgxpool.Pool, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByUsername(pool, input.Username)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Generate token
		token, err := auth.GenerateToken(user.ID, user.Username, jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  user,
		})
	}
}

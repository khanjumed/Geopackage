package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/khanjumed/geopackage/internal/models"
	"github.com/khanjumed/geopackage/internal/services"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashed, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	req.Password = hashed
	if err := models.CreateUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println("🔍 Login attempt for email:", req.Email)
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		fmt.Println("❌ User not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	fmt.Println("✅ User found:", user.Email)

	if !services.CheckPassword(req.Password, user.Password) {
		fmt.Println("❌ Password mismatch for user:", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	fmt.Println("✅ Password matched for:", req.Email)

	token, err := services.GenerateJWT(user.ID, user.Role, os.Getenv("JWT_SECRET"), time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

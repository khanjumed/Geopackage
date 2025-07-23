package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain text password and returns its bcrypt hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("❌ Error while hashing password:", err)
	}
	return string(bytes), err
}

// CheckPassword compares a plaintext password with the hashed one stored in DB
func CheckPassword(input, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	if err != nil {
		fmt.Println("❌ Password does not match:", err)
	} else {
		fmt.Println("✅ Password matched successfully")
	}
	return err == nil
}

// GenerateJWT creates a JWT token with user ID and role
func GenerateJWT(userID int, role string, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("❌ Error generating JWT token:", err)
		return "", err
	}

	fmt.Println("✅ JWT token generated for user:", userID)
	return signed, nil
}

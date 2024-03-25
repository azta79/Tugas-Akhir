package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT_KEY adalah kunci rahasia yang digunakan untuk menandatangani dan memverifikasi token JWT
var JWT_KEY = []byte("983948fdko9f9g03249fo9f9g7")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// Fungsi untuk menghasilkan token JWT
func GenerateToken(userID uint, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expired in 24 hours
	})
	return token.SignedString(JWT_KEY)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
    // Initialize an instance of your custom claims struct
    claims := &Claims{}

    // Parse the token with custom claims
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return JWT_KEY, nil
    })

    return token, err
}

func ExtractToken(authorizationHeader string) (*jwt.Token, error) {
	// Buka header Authorization dan ambil token JWT
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == authorizationHeader {
		return nil, fmt.Errorf("invalid token format")
	}
	

	// Verifikasi token dan kembalikan hasilnya
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
}

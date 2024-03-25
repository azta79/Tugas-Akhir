package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	// Mendapatkan token dari header Authorization
	tokenString := c.GetHeader("Authorization")
	log.Println("Token diperoleh dari header Authorization:", tokenString)

	// Memeriksa keberadaan token
	if tokenString == "" {
		log.Println("Token tidak ditemukan dalam header Authorization")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token diperlukan untuk autentikasi"})
		c.Abort()
		return
	}

	// Menghapus "Bearer " dari token (jika ada)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	log.Println("Token setelah menghapus 'Bearer ':", tokenString)

	// Verifikasi token
	token, err := VerifyToken(tokenString)
	if err != nil {
		log.Println("Gagal memverifikasi token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}

	// Mengambil klaim dari token
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Token tidak valid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}

	log.Println("Token berhasil diverifikasi")

	// Lanjutkan ke penanganan berikutnya jika token valid
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Next()
}

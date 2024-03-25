package controllers

import (
	"MyGram/auth"
	"MyGram/config"
	"MyGram/models"

	"net/http"
	"strconv"
	"time"
	"strings"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {

	// Periksa apakah ini adalah permintaan registrasi
	if c.Request.URL.Path != "/users/register" {
		// Jika bukan permintaan registrasi, kembalikan respons 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	var requestBody struct {
		DOB      string   `json:"dob" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
	}

	// Bind request body ke struct requestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Umur, Email, Password, dan Usernaname harus diisi"})
		return
	}

	// Validasi umur (DOB harus di atas 8 tahun)
    dobTime, err := time.Parse("2006-01-02", requestBody.DOB)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal lahir salah"})
        return
    }
    if time.Since(dobTime).Hours()/24/365 <= 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Umur harus di atas 8 tahun"})
        return
    }

	// Validasi password (minimal 6 karakter)
	if len(requestBody.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal harus 6 karakter"})
		return
	}

	// Hash password sebelum disimpan ke database
	hashedPassword, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		// Handle error
		return
	}

    // Buat objek user baru
    newUser := models.User{
        DOB:       dobTime,
        Email:     requestBody.Email,
        Password:  string(hashedPassword),
        Username:  requestBody.Username,
    }

// Simpan user baru ke dalam database
if err := config.DB.Create(&newUser).Error; err != nil {
    // Periksa apakah kesalahan disebabkan oleh duplikasi email
    if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
        return
    }
    // Periksa apakah kesalahan disebabkan oleh duplikasi username
    if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "username") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
        return
    }
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna baru"})
    return
}


	// Hilangkan field password dari respons
	newUser.Password = ""
	

	// Kirim respons 201 Created dengan data pengguna yang baru dibuat
	userResponse := ConvertToUserResponse(newUser)
    c.JSON(http.StatusCreated, userResponse)
}

func ConvertToUserResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		DOB:      user.DOB,
	}
}



func LoginUser(c *gin.Context) {
	var reqBody models.LoginRequest
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau password belum diisi"})
		return
	}


	// Cari pengguna berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", reqBody.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Periksa kecocokan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate token JWT
	tokenString, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Kirim token JWT sebagai respons
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func GetUser(c *gin.Context) {
	// Mendapatkan ID pengguna dari URL parameter
	userID := strings.TrimPrefix(c.Param("id"), ":")

	// Cari pengguna berdasarkan ID
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hilangkan field password dari respons
	user.Password = ""

	// Buat respons yang diinginkan tanpa deleted_at dan profile_image_url
	userResponse := models.UserResponseGet{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		DOB:       user.DOB,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Kirim informasi pengguna sebagai respons
	c.JSON(http.StatusOK, userResponse)
}


func UpdateUser(c *gin.Context) {
	// Mendapatkan ID pengguna dari URL parameter
	userIDStr := c.Param("id")
userID, err := strconv.Atoi(userIDStr)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    return
}

	// Mendapatkan data request body
	var reqBody models.UpdateUserRequest
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mendapatkan data otorisasi dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Memverifikasi dan mendapatkan userID dari token JWT
	token, err := auth.ExtractToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	authUserIDFloat := claims["user_id"].(float64)
	authUserID := uint(authUserIDFloat)

	

	// Memeriksa apakah userID dalam token sama dengan userID dalam URL parameter
	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update user"})
		return
	}

	// Memperbarui informasi pengguna dalam basis data
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update informasi pengguna
	user.Email = reqBody.Email
	user.Username = reqBody.Username
	user.ProfileImageURL = reqBody.PhotoImageURL
	user.UpdatedAt = time.Now()

	// Simpan perubahan ke dalam basis data
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Kirim respons dengan informasi pengguna yang telah diperbarui
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"username":  user.Username,
		"dob":       user.DOB,
		"update_at": user.UpdatedAt,
	})
}
func DeleteUser(c *gin.Context) {
    userID := c.Param("id") // Menggunakan "id" sebagai parameter
    fmt.Println("User ID from URL:", userID) // Tambahkan log untuk menampilkan ID pengguna dari URL parameter

    // Cari pengguna berdasarkan ID
    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        fmt.Println("Error querying user:", err) // Tambahkan log untuk menampilkan kesalahan saat query
        c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
        return
    }

    // Hapus pengguna dari basis data
    if err := config.DB.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
        return
    }

    // Kirim respons sukses
    c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}





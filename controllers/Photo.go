package controllers

import (
	"MyGram/auth"
	"MyGram/config"
	"MyGram/models"
	"net/http"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


func CreatePhoto(c *gin.Context) {
    var reqBody models.CreatePhotoRequest
    if err := c.BindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Mendapatkan token dari header Authorization
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Butuh token untuk masuk"})
        return
    }

    // Memverifikasi token
    token, err := auth.ExtractToken(authHeader)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token keliru"})
        return
    }

    // Mendapatkan user ID dari token
    claims := token.Claims.(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    // Membuat objek photo baru
    newPhoto := models.Photo{
        Title:     reqBody.Title,
        Caption:   reqBody.Caption,
        PhotoURL:  reqBody.PhotoURL,
        UserID:    userID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(), // Set UpdatedAt menjadi nil saat membuat foto baru
    }

    // Menyimpan photo baru ke dalam database
    if err := config.DB.Create(&newPhoto).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal, pasti ada yang salah"})
        return
    }

    // Mengirim respons 201 Created dengan data photo yang baru dibuat
    c.JSON(http.StatusCreated, newPhoto)
}


func GetPhotos(c *gin.Context) {
    var photos []models.Photo

    // Mengambil semua foto dari database beserta data pengguna terkait
    if err := config.DB.Preload("User").Find(&photos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil foto"})
        return
    }

    // Membuat slice baru untuk menyimpan data foto beserta informasi pengguna yang disaring
    var photosResponse []models.PhotoResponse

    // Memetakan data foto ke dalam response dengan mengambil hanya username dan email dari pengguna
    for _, photo := range photos {
        userBrief := models.UserBriefPhoto{
            Username: photo.User.Username,
            Email:    photo.User.Email,
        }
        photoResponse := models.PhotoResponse{
            ID:        photo.ID,
            Title:     photo.Title,
            Caption:   photo.Caption,
            PhotoURL:  photo.PhotoURL,
            UserID:    photo.UserID,
            User:      userBrief,
            CreatedAt: photo.CreatedAt,
            UpdatedAt: photo.UpdatedAt,
        }
        photosResponse = append(photosResponse, photoResponse)
    }

    // Mengirim data foto beserta informasi pengguna sebagai respons
    c.JSON(http.StatusOK, gin.H{"data": photosResponse})
}





func GetPhoto(c *gin.Context) {
    // Mendapatkan ID photo dari URL parameter
    photoID := c.Param("id")

    // Mengambil photo dari database berdasarkan ID
    var photo models.Photo
    if err := config.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    // Mengirim data photo sebagai respons
    c.JSON(http.StatusOK, gin.H{
        "id":         photo.ID,
        "title":      photo.Title,
        "caption":    photo.Caption,
        "photo_url":  photo.PhotoURL,
        "user_id":    photo.UserID,
        "update_at":  photo.UpdatedAt,
    })
}



func UpdatePhoto(c *gin.Context) {
    // Mendapatkan ID photo dari URL parameter
    photoID := c.Param("id")

    // Mendapatkan data request body
    var reqBody models.UpdatePhotoRequest
    if err := c.BindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Mendapatkan user ID dari konteks autentikasi
    userID, _ := c.Get("user_id")

    // Mengambil photo dari database berdasarkan ID
    var photo models.Photo
    if err := config.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    // Memeriksa apakah pengguna memiliki izin untuk mengupdate photo
    if photo.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update photo"})
        return
    }

    // Update data photo
    photo.Title = reqBody.Title
    photo.Caption = reqBody.Caption
    photo.PhotoURL = reqBody.PhotoURL

    // Menyimpan perubahan ke dalam database
    if err := config.DB.Save(&photo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
        return
    }

    // Mengirim data photo yang telah diupdate sebagai respons
    // tanpa menyertakan informasi pengguna yang kosong
    updatedPhoto := struct {
        ID        uint       `json:"id"`
        Title     string     `json:"title"`
        Caption   string     `json:"caption"`
        PhotoURL  string     `json:"photo_url"`
        UserID    uint       `json:"user_id"`
        UpdateAt  time.Time  `json:"update_at"`
    }{
        ID:        photo.ID,
        Title:     photo.Title,
        Caption:   photo.Caption,
        PhotoURL:  photo.PhotoURL,
        UserID:    photo.UserID,
        UpdateAt:  photo.UpdatedAt,
    }
    c.JSON(http.StatusOK, updatedPhoto)
}


func DeletePhoto(c *gin.Context) {
    // Mendapatkan ID photo dari URL parameter
    photoID := c.Param("photoId")

    // Mendapatkan user ID dari konteks autentikasi
    userID, _ := c.Get("user_id")

    // Mengambil photo dari database berdasarkan ID
    var photo models.Photo
    if err := config.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    // Memeriksa apakah pengguna memiliki izin untuk menghapus photo
    if photo.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete photo"})
        return
    }

    // Menghapus photo dari database
    if err := config.DB.Delete(&photo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
        return
    }

    // Mengirim respons sukses
    c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}



package controllers

import (
	"MyGram/config"
	"MyGram/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateSocialMedia creates a new social media entry

func CreateSocialMedia(c *gin.Context) {
    var requestData models.SocialMedia
    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Retrieve user ID from JWT token
    userID, _ := c.Get("user_id")
    requestData.UserID = userID.(uint)

    // Create social media entry in database
    if err := config.DB.Create(&requestData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media"})
        return
    }

    // Create simplified social media response
    response := models.SocialMediaResponse{
        ID:              requestData.ID,
        Name:            requestData.Name,
        SocialMediaURL:  requestData.SocialMediaURL,
        UserID:          requestData.UserID,
        CreatedAt:       requestData.CreatedAt,
    }

    // Convert response to JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
        return
    }

    // Send JSON response
    c.Data(http.StatusCreated, "application/json", jsonResponse)
}

type SocialMediaResponseGet struct {
    SocialMedia []models.SocialMedia `json:"social-media"`
}


// GetSocialMedias retrieves all social media entries
func GetSocialMedias(c *gin.Context) {
    var socialMedias []models.SocialMedia
    if err := config.DB.Preload("User").Find(&socialMedias).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social medias"})
        return
    }

    // Create the response structure
    response := struct {
        SocialMedia []struct {
            ID             uint      `json:"id"`
            Name           string    `json:"name"`
            SocialMediaURL string    `json:"social_media_url"`
            UserID         uint      `json:"user_id"`
            CreatedAt      time.Time `json:"created_at"`
            UpdatedAt      time.Time `json:"updated_at"`
            User           struct {
                ID              uint   `json:"id"`
                Username        string `json:"username"`
                ProfileImageURL string `json:"profile_image_url"`
            } `json:"user"`
        } `json:"social-media"`
    }{}

    for _, sm := range socialMedias {
        response.SocialMedia = append(response.SocialMedia, struct {
            ID             uint      `json:"id"`
            Name           string    `json:"name"`
            SocialMediaURL string    `json:"social_media_url"`
            UserID         uint      `json:"user_id"`
            CreatedAt      time.Time `json:"created_at"`
            UpdatedAt      time.Time `json:"updated_at"`
            User           struct {
                ID              uint   `json:"id"`
                Username        string `json:"username"`
                ProfileImageURL string `json:"profile_image_url"`
            } `json:"user"`
        }{
            ID:             sm.ID,
            Name:           sm.Name,
            SocialMediaURL: sm.SocialMediaURL,
            UserID:         sm.UserID,
            CreatedAt:      sm.CreatedAt,
            UpdatedAt:      sm.UpdatedAt,
            User: struct {
                ID              uint   `json:"id"`
                Username        string `json:"username"`
                ProfileImageURL string `json:"profile_image_url"`
            }{
                ID:              sm.User.ID,
                Username:        sm.User.Username,
                ProfileImageURL: sm.User.ProfileImageURL,
            },
        })
    }

    c.JSON(http.StatusOK, response)
}





// GetSocialMedia retrieves a specific social media entry by ID
func GetSocialMedia(c *gin.Context) {
    socialMediaID := c.Param("id")
    var socialMedia models.SocialMedia
    if err := config.DB.Preload("User").Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
        return
    }

     // Customize the user response within social media
     userResponse := struct {
        ID             uint   `json:"id"`
        Username       string `json:"username"`
        ProfileImageURL string `json:"profile_image_url"`
    }{
        ID:             socialMedia.User.ID,
        Username:       socialMedia.User.Username,
        ProfileImageURL: socialMedia.User.ProfileImageURL,
    }

    // Customize the social media response
    socialMediaResponse := struct {
        ID             uint      `json:"id"`
        Name           string    `json:"name"`
        SocialMediaURL string    `json:"social_media_url"`
        UserID         uint      `json:"user_id"`
        CreatedAt      time.Time `json:"created_at"`
        UpdatedAt      time.Time `json:"updated_at"`
        User           interface{} `json:"user"`
    }{
        ID:             socialMedia.ID,
        Name:           socialMedia.Name,
        SocialMediaURL: socialMedia.SocialMediaURL,
        UserID:         socialMedia.UserID,
        CreatedAt:      socialMedia.CreatedAt,
        UpdatedAt:      socialMedia.UpdatedAt,
        User:           userResponse,
    }

    // Memperbarui data pengguna untuk menyertakan URL gambar profil
    socialMedia.User.ProfileImageURL = "https://example.com/profile.jpg" // Ganti dengan URL gambar profil yang sesuai

    // Mengirimkan respons dengan data social media yang telah dimuat
    c.JSON(http.StatusOK, gin.H{"social-media": []interface{}{socialMediaResponse}})
}



// UpdateSocialMedia updates an existing social media entry
func UpdateSocialMedia(c *gin.Context) {
    socialMediaID := c.Param("id")
    var updateData models.SocialMedia
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Retrieve user ID from JWT token
    userID, _ := c.Get("user_id")

    // Check if the social media entry exists
    var existingSocialMedia models.SocialMedia
    if err := config.DB.Where("id = ?", socialMediaID).First(&existingSocialMedia).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
        return
    }

    // Check if the user is authorized to update this social media entry
    if existingSocialMedia.UserID != userID.(uint) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this social media"})
        return
    }

    // Update social media entry in database
    if err := config.DB.Model(&existingSocialMedia).Updates(&updateData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media"})
        return
    }

    // Create SocialMediaResponse
    response := models.SocialMediaResponseUpdate{
        ID:             existingSocialMedia.ID,
        Name:           updateData.Name,
        SocialMediaURL: updateData.SocialMediaURL,
        UserID:         existingSocialMedia.UserID,
        UpdatedAt:      existingSocialMedia.UpdatedAt.Format(time.RFC3339), // Format the time as needed
    }

    c.JSON(http.StatusOK, response)
}


// DeleteSocialMedia deletes a social media entry
func DeleteSocialMedia(c *gin.Context) {
    socialMediaID := c.Param("id")

    // Retrieve user ID from JWT token
    userID, _ := c.Get("user_id")

    // Check if the social media entry exists
    var existingSocialMedia models.SocialMedia
    if err := config.DB.Where("id = ?", socialMediaID).First(&existingSocialMedia).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
        return
    }

    // Check if the user is authorized to delete this social media entry
    if existingSocialMedia.UserID != userID.(uint) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this social media"})
        return
    }
}
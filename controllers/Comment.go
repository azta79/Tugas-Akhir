package controllers

import (
	"MyGram/auth"
	"MyGram/config"
	"MyGram/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
)



func CreateComment(c *gin.Context) {
    // Menjalankan middleware autentikasi JWT
    auth.AuthMiddleware(c)

    // Mendapatkan data request body
    var reqBody models.CreateCommentRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Mendapatkan ID pengguna dari konteks
    userID, _ := c.Get("user_id")
   // fmt.Println("User ID:", userID)

    // Membuat komentar baru
    comment := models.Comment{
        UserID:    userID.(uint),
        PhotoID:   reqBody.PhotoID,
        Message:   reqBody.Message,
        CreatedAt: time.Now(),
    }
   
    // Menyimpan komentar ke database
    if err := config.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
        return
    }
    
    // Mengonversi ke CreateCommentResponse
    createCommentResponse := models.CreateCommentResponse{
        ID:        comment.ID,
        UserID:    comment.UserID,
        PhotoID:   comment.PhotoID,
        Message:   comment.Message,
        CreatedAt: comment.CreatedAt,
    }

    // Mengirim respons JSON dengan komentar yang telah dibuat
    c.JSON(http.StatusCreated, createCommentResponse)
}



// GetComments retrieves all comments with associated user and photo
func GetComments(c *gin.Context) {
        // Menjalankan middleware autentikasi JWT
        auth.AuthMiddleware(c)
    var comments []models.Comment
    if err := config.DB.Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
        return
    }

    // Membuat slice untuk menyimpan respons komentar
    var commentResponses []models.CommentResponse

    // Iterasi melalui setiap komentar dan membuat responsnya
    for _, comment := range comments {
        // Dapatkan data pengguna terkait
        var user models.User
        if err := config.DB.First(&user, comment.UserID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated user"})
            return
        }

        // Dapatkan data foto terkait
        var photo models.Photo
        if err := config.DB.First(&photo, comment.PhotoID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated photo"})
            return
        }

        // Buat respons komentar dan tambahkan ke slice
        commentResponse := models.CommentResponse{
            ID:        comment.ID,
            UserID:    comment.UserID,
            PhotoID:   comment.PhotoID,
            Message:   comment.Message,
            CreatedAt: comment.CreatedAt,
            UpdatedAt: comment.UpdatedAt,
            User: models.UserBriefComment{
                ID:       user.ID,
                Username: user.Username,
                Email:    user.Email,
            },
            Photo: models.PhotoResponseComment{
                ID:        photo.ID,
                Title:     photo.Title,
                Caption:   photo.Caption,
                PhotoURL:  photo.PhotoURL,
                UserID:    photo.UserID,
            },
        }
        commentResponses = append(commentResponses, commentResponse)
    }

    // Kirim respons dengan semua komentar
    c.JSON(http.StatusOK, commentResponses)
}




// GetComment retrieves a specific comment with associated user and photo
func GetComment(c *gin.Context) {
        // Menjalankan middleware autentikasi JWT
        auth.AuthMiddleware(c)
    // Mendapatkan ID komentar dari URL parameter
    commentID := c.Param("id")

    // Cari komentar berdasarkan ID
    var comment models.Comment
    if err := config.DB.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

    // Dapatkan data pengguna terkait
    var user models.User
    if err := config.DB.First(&user, comment.UserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated user"})
        return
    }

    // Dapatkan data foto terkait
    var photo models.Photo
    if err := config.DB.First(&photo, comment.PhotoID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated photo"})
        return
    }

    // Membuat respons sesuai dengan spesifikasi
    response := models.CommentResponse{
        ID:        comment.ID,
        UserID:    comment.UserID,
        PhotoID:   comment.PhotoID,
        Message:   comment.Message,
        CreatedAt: comment.CreatedAt,
        UpdatedAt: comment.UpdatedAt,
        User: models.UserBriefComment{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
        },
        Photo: models.PhotoResponseComment{
            ID:        photo.ID,
            Title:     photo.Title,
            Caption:   photo.Caption,
            PhotoURL:  photo.PhotoURL,
            UserID:    photo.UserID,
        },
    }

    // Kirim respons
    c.JSON(http.StatusOK, response)
}


func UpdateComment(c *gin.Context) {
    // Menjalankan middleware autentikasi JWT
    auth.AuthMiddleware(c)

    commentID := c.Param("id")
    var reqBody models.UpdateCommentRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var comment models.Comment
    if err := config.DB.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

    // Check if the user is authorized to update the comment
    userID, _ := c.Get("user_id")
    if comment.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update comment"})
        return
    }

   // Update comment message
comment.Message = reqBody.Message

// Simpan perubahan komentar
if err := config.DB.Save(&comment).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
    return
}

// Dapatkan informasi foto terkait dari komentar
var photo models.Photo
if err := config.DB.Model(&comment).Related(&photo, "Photo").Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated photo"})
    return
}

// Kirim respons sukses dengan informasi yang diminta
c.JSON(http.StatusOK, gin.H{
    "id":        comment.ID,
    "title":     photo.Title,
    "message":   comment.Message,
    "caption":   photo.Caption,
    "photo_url": photo.PhotoURL,
    "user_id":   photo.UserID,
    "update_at": comment.UpdatedAt,
})
}


// DeleteComment deletes a specific comment by ID
func DeleteComment(c *gin.Context) {
	    // Menjalankan middleware autentikasi JWT
		auth.AuthMiddleware(c)
    commentID := c.Param("id")
    var comment models.Comment
    if err := config.DB.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

    // Check if the user is authorized to delete the comment
    userID, _ := c.Get("user_id")
    if comment.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete comment"})
        return
    }

    if err := config.DB.Delete(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

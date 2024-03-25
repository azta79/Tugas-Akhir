package models

import "time"


type CommentUser struct {
    ID uint `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
    }

type Comment struct {
    ID        uint       `gorm:"primary_key" json:"id"`
    UserID    uint       `gorm:"not null;foreignkey:UserID" json:"user_id"`
    PhotoID   uint       `gorm:"not null;foreignkey:PhotoID" json:"photo_id"`
    Message   string     `gorm:"not null" json:"message" validate:"required"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `gorm:"deleted_at" json:"-"`
    User      User ` json:"user"`
    Photo     Photo    `json:"photo"`
}

type PhotoComment struct {
    ID        uint       `json:"id"`
    Title     string     `json:"title"`
    Caption   string     `json:"caption"`
    PhotoURL  string     `json:"photo_url"`
    UserID    uint       `json:"user_id"`
}



type CommentResponse struct {
    ID        uint       `json:"id"`
    UserID    uint       `json:"user_id"`
    PhotoID   uint       `json:"photo_id"`
    Message   string     `json:"message"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    User      UserBriefComment        `json:"user"`
    Photo     PhotoResponseComment       `json:"photo"`
}

type CreateCommentRequest struct {
    Message string `json:"message" binding:"required"`
    PhotoID uint   `json:"photo_id" binding:"required"`
}

type UpdateCommentRequest struct {
    Message string `json:"message" binding:"required"`
    
}

type CreateCommentResponse struct {
    ID        uint       `json:"id"`
    UserID    uint       `json:"user_id"`
    PhotoID   uint       `json:"photo_id"`
    Message   string     `json:"message"`
    CreatedAt time.Time  `json:"created_at"`
}

type UserBriefComment struct {
    ID        uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

type PhotoResponseComment struct {
    ID        uint           `json:"id"`
    Title     string         `json:"title"`
    Caption   string         `json:"caption"`
    PhotoURL  string         `json:"photo_url"`
    UserID    uint           `json:"user_id"`
    } 



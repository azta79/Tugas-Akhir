package models

import (
    "time"
    

)

type Photo struct {
    ID        uint       `gorm:"primary_key" json:"id"`
    Title     string     `gorm:"not null" json:"title" validate:"required"`
    Caption   string     `gorm:"not null" json:"caption" validate:"required"`
    PhotoURL  string     `gorm:"not null" json:"photo_url" validate:"required"`
    UserID    uint       `gorm:"not null" json:"user_id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"update_at"`
    DeletedAt *time.Time `gorm:"-" json:"-"`
    User      User  ` json:"user"`
    
}


type UserBriefPhoto struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}


type CreatePhotoRequest struct {
    Title    string `json:"title" binding:"required"`
    Caption  string `json:"caption" binding:"required"`
    PhotoURL string `json:"photo_url" binding:"required"`
}

type UpdatePhotoRequest struct {
    Title    string `json:"title"`
    Caption  string `json:"caption"`
    PhotoURL string `json:"photo_url"`
}

// PhotoResponse adalah struktur untuk menanggapi data foto
type PhotoResponse struct {
    ID        uint           `json:"id"`
    Title     string         `json:"title"`
    Caption   string         `json:"caption"`
    PhotoURL  string         `json:"photo_url"`
    UserID    uint           `json:"user_id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time     ` json:"update_at"`
    User      UserBriefPhoto `json:"user"`
    }

 
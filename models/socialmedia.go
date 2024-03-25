package models

import "time"

type SocialMedia struct {
    ID              uint       `gorm:"primary_key" json:"id"`
    Name            string     `gorm:"not null" json:"name" validate:"required"`
    SocialMediaURL  string     `gorm:"not null" json:"social_media_url" validate:"required"`
    UserID          uint       `gorm:"not null;foreignkey:UserID" json:"user_id"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    DeletedAt *time.Time `gorm:"deleted_at" json:"-"`

        // Bidang ini akan merujuk pada model User
        User         User       `gorm:"foreignkey:UserID" json:"user"`
}

type SocialMediaResponse struct {
    ID              uint      `json:"id"`
    Name            string    `json:"name"`
    SocialMediaURL  string    `json:"social_media_url"`
    UserID           uint   `json:"user_id"`
    CreatedAt       time.Time `json:"created_at"`
}

type SocialMediaResponseUpdate struct {
    ID             uint   `json:"id"`
    Name           string `json:"name"`
    SocialMediaURL string `json:"social_media_url"`
    UserID         uint   `json:"user_id"`
    UpdatedAt      string `json:"updated_at"`
}
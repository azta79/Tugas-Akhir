package models

import (
    "time"
    
    
)

type User struct {
    ID        uint       `json:"id" gorm:"primary_key;column:id"`
    Username  string     `json:"username" gorm:"type:varchar(100);unique;not null" validate:"required,min=3,max=100"`
    Email     string     `json:"email" gorm:"type:varchar(100);unique;not null" validate:"required,email"`
    Password  string     `json:"-" gorm:"type:varchar(100);not null" validate:"required,min=6"`
    DOB       time.Time  `json:"dob" gorm:"type:date"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
    ProfileImageURL string     `json:"profile_image_url"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    DOB       time.Time `json:"dob"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponseGet struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	DOB      time.Time `json:"dob"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type UpdateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
    PhotoImageURL   string `json:"photo_image_url"`
}
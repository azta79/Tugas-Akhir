package repository

import (
	"MyGram/models" // Ganti dengan path package model Anda
	"gorm.io/gorm"
)

// UserRepository adalah repository untuk operasi CRUD pada pengguna
type UserRepository struct {
	DB *gorm.DB
}

// CreateUser digunakan untuk membuat pengguna baru dalam database
func (ur *UserRepository) CreateUser(user *models.User) error {
	return ur.DB.Create(user).Error
}

// GetUserByID digunakan untuk mendapatkan pengguna berdasarkan ID
func (ur *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := ur.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

package config

import (
    
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "MyGram/models"
)

var (
    DB *gorm.DB
)

func InitDB() {
    var err error
    DB, err = gorm.Open("mysql", "WindowsX:12345@tcp(localhost:3306)/db_MyGram?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        panic("failed to connect database")
    }

    // Auto migrate models
    DB.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

    }

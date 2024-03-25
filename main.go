package main

import (
    "MyGram/config"
    "MyGram/router"
)

func main() {
    // Initialize database
    config.InitDB()

    // Initialize Gin router
    r := router.Init()

    // Run server
    r.Run(":8888")
}
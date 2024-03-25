package router

import (
    "MyGram/controllers"
    "MyGram/auth"

    "github.com/gin-gonic/gin"
)

// Init initializes the Gin router
func Init() *gin.Engine {
    // Initialize Gin router
    r := gin.Default()

// Routes
r.POST("/users/register", controllers.RegisterUser)
r.POST("/users/login", controllers.LoginUser)

authorized := r.Group("/")
{
	authorized.Use(auth.AuthMiddleware) // Middleware autentikasi JWT

	authorized.GET("/users/:id", controllers.GetUser)
	authorized.PUT("/users/:id", controllers.UpdateUser)
	authorized.DELETE("/users/:id", controllers.DeleteUser)

	authorized.POST("/photos", controllers.CreatePhoto)
	authorized.GET("/photos", controllers.GetPhotos)
	authorized.GET("/photos/:id", controllers.GetPhoto)
	authorized.PUT("/photos/:id", controllers.UpdatePhoto)
	authorized.DELETE("/photos/:id", controllers.DeletePhoto)

	authorized.POST("/comments", controllers.CreateComment)
	authorized.GET("/comments", controllers.GetComments)
	authorized.GET("/comments/:id", controllers.GetComment)
	authorized.PUT("/comments/:id", controllers.UpdateComment)
	authorized.DELETE("/comments/:commentId", controllers.DeleteComment)

	authorized.POST("/social-medias", controllers.CreateSocialMedia)
	authorized.GET("/social-medias", controllers.GetSocialMedias)
	authorized.GET("/social-medias/:id", controllers.GetSocialMedia)
	authorized.PUT("/social-medias/:id", controllers.UpdateSocialMedia)
	authorized.DELETE("/social-medias/:id", controllers.DeleteSocialMedia)
}

return r

}

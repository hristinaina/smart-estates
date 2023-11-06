package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"smarthome-back/controllers"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.GET("/:id", userController.GetUser)
	}
}


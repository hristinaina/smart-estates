package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"smarthome-back/controllers"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, database *sql.DB) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db, database)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.GET("/test", userController.TestGetMethod)
	}

	realEstateRoutes := r.Group("/api/realestate")
	{
		realEstateController := controllers.NewRealEstateController(db, database)
		realEstateRoutes.GET("/:userId", realEstateController.GetAll)
		realEstateRoutes.GET("/:id", realEstateController.Get)
	}
}

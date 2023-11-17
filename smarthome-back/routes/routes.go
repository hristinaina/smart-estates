package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"smarthome-back/controllers"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.GET("/test", userController.TestGetMethod)
	}

	realEstateRoutes := r.Group("/api/real-estates")
	{
		realEstateController := controllers.NewRealEstateController(db)
		realEstateRoutes.GET("/", realEstateController.GetAll)
		realEstateRoutes.GET("/user/:userId", realEstateController.GetAllByUserId)
		realEstateRoutes.GET("/:id", realEstateController.Get)
		realEstateRoutes.PUT("/:id/:state", realEstateController.ChangeState)
		realEstateRoutes.POST("/", realEstateController.Add)
	}

	deviceRoutes := r.Group("/api/devices")
	{
		deviceController := controllers.NewDeviceController(db)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/estate/:estateId", deviceController.GetAllByEstateId)
	}
}

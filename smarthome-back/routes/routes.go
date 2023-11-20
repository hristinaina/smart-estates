package routes

import (
	"database/sql"
	"smarthome-back/controllers"
	"smarthome-back/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.GET("/test", userController.TestGetMethod)

		authController := controllers.NewAuthController(db)
		middleware := middleware.NewMiddleware(db)
		userRoutes.POST("/login", authController.Login)
		userRoutes.GET("/validate", middleware.RequireAuth, authController.Validate)
		userRoutes.POST("/logout", authController.Logout)
		userRoutes.POST("/verificationMail", authController.SendVerificationMail)
		userRoutes.POST("/activate", authController.ActivateAccount)
	}

	realEstateRoutes := r.Group("/api/real-estates")
	{
		realEstateController := controllers.NewRealEstateController(db)
		realEstateRoutes.GET("/", realEstateController.GetAll)
		realEstateRoutes.GET("/user/:userId", realEstateController.GetAllByUserId)
		realEstateRoutes.GET("/:id", realEstateController.Get)
		realEstateRoutes.PUT("/:id/:state", middleware.AdminMiddleware, realEstateController.ChangeState) // user can't use this
		realEstateRoutes.POST("/", middleware.UserMiddleware, realEstateController.Add)                   // admin can't use this
	}

	deviceRoutes := r.Group("/api/devices")
	{
		deviceController := controllers.NewDeviceController(db)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/estate/:estateId", deviceController.GetAllByEstateId)
		deviceRoutes.POST("/", deviceController.Add)
	}
	airConditionerRoutes := r.Group("/api/ac")
	{
		airConditionerController := controllers.NewAirConditionerController(db)
		airConditionerRoutes.GET("/:id", airConditionerController.Get)
	}
}

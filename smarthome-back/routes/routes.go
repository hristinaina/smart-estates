package routes

import (
	"database/sql"
	"smarthome-back/controllers"
	devicesController "smarthome-back/controllers/devices"
	"smarthome-back/middleware"
	"smarthome-back/mqtt_client"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, mqtt *mqtt_client.MQTTClient) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.GET("/test", userController.TestGetMethod)
		userRoutes.POST("/verify-email", userController.SendResetPasswordEmail)
		userRoutes.POST("/reset-password", userController.ResetPassword)

		// todo promeni middleware
		authController := controllers.NewAuthController(db)
		middleware := middleware.NewMiddleware(db)
		userRoutes.POST("/login", authController.Login)
		userRoutes.GET("/validate", middleware.RequireAuth, authController.Validate)
		userRoutes.POST("/logout", middleware.RequireAuth, authController.Logout)
		userRoutes.POST("/verificationMail", authController.SendVerificationMail)
		userRoutes.POST("/activate", authController.ActivateAccount)

		superadminController := controllers.NewSuperAdminController(db)
		userRoutes.POST("/reset-superadmin-password", middleware.SuperAdminMiddleware, superadminController.ResetPassword)
		userRoutes.POST("/add-admin", middleware.SuperAdminMiddleware, superadminController.AddAdmin)
		userRoutes.POST("/edit-admin", middleware.SuperAdminMiddleware, superadminController.EditSuperAdmin)
	}

	realEstateRoutes := r.Group("/api/real-estates")
	{
		realEstateController := controllers.NewRealEstateController(db)
		realEstateRoutes.GET("/", realEstateController.GetAll)
		realEstateRoutes.GET("/user/:userId", realEstateController.GetAllByUserId)
		realEstateRoutes.GET("/:id", realEstateController.Get)
		realEstateRoutes.GET("/pending", realEstateController.GetPending)
		realEstateRoutes.PUT("/:id/:state", realEstateController.ChangeState) // user can't use this
		realEstateRoutes.POST("/", realEstateController.Add)                  // admin can't use this
	}

	deviceRoutes := r.Group("/api/devices")
	{
		deviceController := controllers.NewDeviceController(db, mqtt)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/", deviceController.GetAll)
		deviceRoutes.GET("/estate/:estateId", deviceController.GetAllByEstateId)
		deviceRoutes.POST("/", deviceController.Add)
	}
	airConditionerRoutes := r.Group("/api/ac")
	{
		airConditionerController := controllers.NewAirConditionerController(db)
		airConditionerRoutes.GET("/:id", airConditionerController.Get)
	}

	uploadImageRoutes := r.Group("/api/upload")
	{
		imageUploadController := controllers.NewImageController()
		uploadImageRoutes.POST("/:real-estate-name", imageUploadController.Post)
		uploadImageRoutes.GET("/:file-name", imageUploadController.Get)
	}

	lampRoutes := r.Group("api/lamp")
	{
		lampController := devicesController.NewLampController(db)
		lampRoutes.GET("/:id", lampController.Get)
		lampRoutes.GET("/", lampController.GetAll)
		lampRoutes.PUT("/on/:id", lampController.TurnOn)
		lampRoutes.PUT("/off/:id", lampController.TurnOff)
		lampRoutes.PUT("/:id/:level", lampController.SetLightning)
		lampRoutes.POST("/", lampController.Add)
	}
}

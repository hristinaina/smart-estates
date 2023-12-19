package routes

import (
	"database/sql"
	"smarthome-back/controllers"
	"smarthome-back/middleware"
	"smarthome-back/mqtt_client"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) {
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
		deviceController := controllers.NewDeviceController(db, mqtt, influxDb)
		middleware := middleware.NewMiddleware(db)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/", deviceController.GetAll)
		deviceRoutes.GET("/estate/:estateId", middleware.RequireAuth, deviceController.GetAllByEstateId)
		deviceRoutes.POST("/", middleware.RequireAuth, deviceController.Add)
	}
	airConditionerRoutes := r.Group("/api/ac")
	{
		airConditionerController := controllers.NewAirConditionerController(db, mqtt)
		airConditionerRoutes.GET("/:id", airConditionerController.Get)
		airConditionerRoutes.PUT("history", airConditionerController.GetHistoryData)
	}
	solarPanelRoutes := r.Group("/api/sp")
	{
		middleware := middleware.NewMiddleware(db)
		SolarPanelController := controllers.NewSolarPanelController(db, influxDb)
		solarPanelRoutes.GET("/:id", SolarPanelController.Get)
		solarPanelRoutes.PUT("/graphData", middleware.RequireAuth, SolarPanelController.GetGraphData)
	}
	uploadImageRoutes := r.Group("/api/upload")
	{
		imageUploadController := controllers.NewImageController()
		uploadImageRoutes.POST("/:real-estate-name", imageUploadController.Post)
		uploadImageRoutes.GET("/:file-name", imageUploadController.Get)
	}

	ambientSensor := r.Group("/api/ambient")
	{
		ambientSensorController := controllers.NewAmbientSensorController(db, mqtt)
		ambientSensor.GET("/:id", ambientSensorController.Get)
		ambientSensor.GET("/last-hour/:id", ambientSensorController.GetValueForHour)
		ambientSensor.POST("/selected-time/:id", ambientSensorController.GetValueForSelectedTime)
		ambientSensor.POST("/selected-date/:id", ambientSensorController.GetValuesForDate)
	}
}

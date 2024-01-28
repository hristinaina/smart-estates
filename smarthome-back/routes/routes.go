package routes

import (
	"database/sql"
	"smarthome-back/controllers"
	devicesController "smarthome-back/controllers/devices"
	"smarthome-back/middleware"
	"smarthome-back/mqtt_client"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db)
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
		deviceController := devicesController.NewDeviceController(db, mqtt, influxDb)
		middleware := middleware.NewMiddleware(db)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/", deviceController.GetAll)
		deviceRoutes.GET("/estate/:estateId", middleware.RequireAuth, deviceController.GetAllByEstateId)
		deviceRoutes.POST("/", middleware.RequireAuth, deviceController.Add)
		deviceRoutes.GET("/consumption-device/:id", deviceController.GetConsumptionDeviceDto)
	}
	airConditionerRoutes := r.Group("/api/ac")
	{
		airConditionerController := devicesController.NewAirConditionerController(db, mqtt)
		airConditionerRoutes.GET("/:id", airConditionerController.Get)
		airConditionerRoutes.PUT("history", airConditionerController.GetHistoryData)
	}
	solarPanelRoutes := r.Group("/api/sp")
	{
		middleware := middleware.NewMiddleware(db)
		SolarPanelController := devicesController.NewSolarPanelController(db, influxDb)
		solarPanelRoutes.GET("/:id", SolarPanelController.Get)
		solarPanelRoutes.PUT("/graphData", middleware.RequireAuth, SolarPanelController.GetGraphData)
		solarPanelRoutes.GET("/lastValue/:id", middleware.RequireAuth, SolarPanelController.GetValueFromLastMinute)
	}
	homeBatteryRoutes := r.Group("/api/hb")
	{
		middleware := middleware.NewMiddleware(db)
		HomeBatteryController := devicesController.NewHomeBatteryController(db, influxDb)
		homeBatteryRoutes.GET("/:id", middleware.RequireAuth, HomeBatteryController.Get)
		homeBatteryRoutes.GET("/last-hour/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForLastHour)
		homeBatteryRoutes.POST("/selected-time/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForSelectedTime)
		homeBatteryRoutes.POST("/selected-date/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForSelectedDate)
	}
	uploadImageRoutes := r.Group("/api/upload")
	{
		imageUploadController := controllers.NewImageController()
		uploadImageRoutes.POST("/:real-estate-name", imageUploadController.Post)
		uploadImageRoutes.GET("/:file-name", imageUploadController.Get)
	}

	ambientSensor := r.Group("/api/ambient")
	{
		ambientSensorController := devicesController.NewAmbientSensorController(db, mqtt)
		ambientSensor.GET("/:id", ambientSensorController.Get)
		ambientSensor.GET("/last-hour/:id", ambientSensorController.GetValueForHour)
		ambientSensor.POST("/selected-time/:id", ambientSensorController.GetValueForSelectedTime)
		ambientSensor.POST("/selected-date/:id", ambientSensorController.GetValuesForDate)
	}

	lampRoutes := r.Group("api/lamp")
	{
		lampController := devicesController.NewLampController(db, influxDb)
		lampRoutes.GET("/:id", lampController.Get)
		lampRoutes.GET("/", lampController.GetAll)
		lampRoutes.PUT("/on/:id", lampController.TurnOn)
		lampRoutes.PUT("/off/:id", lampController.TurnOff)
		lampRoutes.PUT("/:id/:level", lampController.SetLightning)
		lampRoutes.POST("/", lampController.Add)
		lampRoutes.DELETE("/:id", lampController.Delete)
		lampRoutes.GET("/graph/:id/:from/:to", lampController.GetGraphData)
	}

	vehicleGateRoutes := r.Group("api/vehicle-gate")
	{
		vehicleGateController := devicesController.NewVehicleGateController(db, influxDb)
		vehicleGateRoutes.GET("/:id", vehicleGateController.Get)
		vehicleGateRoutes.GET("/", vehicleGateController.GetAll)
		vehicleGateRoutes.PUT("/open/:id", vehicleGateController.Open)
		vehicleGateRoutes.PUT("/close/:id", vehicleGateController.Close)
		vehicleGateRoutes.PUT("/private/:id", vehicleGateController.ToPrivate)
		vehicleGateRoutes.PUT("/public/:id", vehicleGateController.ToPublic)
		vehicleGateRoutes.POST("/", vehicleGateController.Add)
		vehicleGateRoutes.DELETE("/:id", vehicleGateController.Delete)
		vehicleGateRoutes.GET("/license-plate/:id", vehicleGateController.GetLicensePlates)
		vehicleGateRoutes.POST("/license-plate", vehicleGateController.AddLicensePlate)
		vehicleGateRoutes.GET("/license-plate", vehicleGateController.GetAllLicensePlates)
		vehicleGateRoutes.GET("/count/:id/:from/:to/:license-plate", vehicleGateController.GetLicencePlatesCount)
	}

	washingMachineRoutes := r.Group("/api/wm")
	{
		middleware := middleware.NewMiddleware(db)
		washingMachineController := devicesController.NewWashingMachineController(db, mqtt)
		washingMachineRoutes.GET("/:id", washingMachineController.Get)
		washingMachineRoutes.POST("/schedule", middleware.RequireAuth, washingMachineController.AddScheduledMode)
		washingMachineRoutes.GET("/schedule/:id", middleware.RequireAuth, washingMachineController.GetScheduledModes)
	}
}

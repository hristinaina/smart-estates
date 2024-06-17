package routes

import (
	"database/sql"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	devicesController "smarthome-back/controllers/devices"
	"smarthome-back/middleware"
	"smarthome-back/mqtt_client"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client, cacheService cache.CacheService) {
	userRoutes := r.Group("/api/users")
	{
		userController := controllers.NewUserController(db, cacheService)
		userRoutes.POST("/verify-email", userController.SendResetPasswordEmail)
		userRoutes.POST("/reset-password", userController.ResetPassword)

		authController := controllers.NewAuthController(db, cacheService)
		middleware := middleware.NewMiddleware(db, cacheService)
		userRoutes.POST("/login", authController.Login)
		userRoutes.GET("/validate", middleware.RequireAuth, authController.Validate)
		userRoutes.POST("/logout", middleware.RequireAuth, authController.Logout)
		userRoutes.POST("/verificationMail", authController.SendVerificationMail)
		userRoutes.POST("/activate", authController.ActivateAccount)

		superadminController := controllers.NewSuperAdminController(db, cacheService)
		userRoutes.POST("/reset-superadmin-password", middleware.SuperAdminMiddleware, superadminController.ResetPassword)
		userRoutes.POST("/add-admin", middleware.SuperAdminMiddleware, superadminController.AddAdmin)
		userRoutes.POST("/edit-admin", middleware.SuperAdminMiddleware, superadminController.EditSuperAdmin)
	}

	realEstateRoutes := r.Group("/api/real-estates")
	{
		realEstateController := controllers.NewRealEstateController(db, &cacheService)
		realEstateRoutes.GET("/", realEstateController.GetAll)
		realEstateRoutes.GET("/cities", realEstateController.GetAllCities)
		realEstateRoutes.GET("/user/:userId", realEstateController.GetAllByUserId)
		realEstateRoutes.GET("/:id", realEstateController.Get)
		realEstateRoutes.GET("/pending", realEstateController.GetPending)
		realEstateRoutes.PUT("/:id/:state", realEstateController.ChangeState) // user can't use this
		realEstateRoutes.POST("/", realEstateController.Add)                  // admin can't use this
	}

	deviceRoutes := r.Group("/api/devices")
	{
		deviceController := devicesController.NewDeviceController(db, mqtt, influxDb, cacheService)
		middleware := middleware.NewMiddleware(db, cacheService)
		deviceRoutes.GET("/:id", deviceController.Get)
		deviceRoutes.GET("/", deviceController.GetAll)
		deviceRoutes.GET("/estate/:estateId", middleware.RequireAuth, deviceController.GetAllByEstateId)
		deviceRoutes.POST("/", middleware.RequireAuth, deviceController.Add)
		deviceRoutes.GET("/consumption-device/:id", deviceController.GetConsumptionDeviceDto)
	}
	airConditionerRoutes := r.Group("/api/ac")
	{
		airConditionerController := devicesController.NewAirConditionerController(db, mqtt, cacheService)
		middleware := middleware.NewMiddleware(db, cacheService)
		airConditionerRoutes.GET("/:id", airConditionerController.Get)
		airConditionerRoutes.PUT("history", middleware.RequireAuth, airConditionerController.GetHistoryData)
		airConditionerRoutes.POST("/edit/:id", middleware.RequireAuth, airConditionerController.EditSpecialModes)
	}
	solarPanelRoutes := r.Group("/api/sp")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		SolarPanelController := devicesController.NewSolarPanelController(db, influxDb, cacheService)
		solarPanelRoutes.GET("/:id", SolarPanelController.Get)
		solarPanelRoutes.PUT("/graphData", middleware.RequireAuth, SolarPanelController.GetGraphData)
		solarPanelRoutes.GET("/lastValue/:id", middleware.RequireAuth, SolarPanelController.GetValueFromLastMinute)
		solarPanelRoutes.PUT("/production", middleware.RequireAuth, SolarPanelController.GetProductionForSP)
	}
	homeBatteryRoutes := r.Group("/api/hb")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		HomeBatteryController := devicesController.NewHomeBatteryController(db, influxDb, cacheService)
		homeBatteryRoutes.GET("/:id", middleware.RequireAuth, HomeBatteryController.Get)
		homeBatteryRoutes.GET("/last-hour/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForLastHour)
		homeBatteryRoutes.POST("/selected-time/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForSelectedTime)
		homeBatteryRoutes.POST("/selected-date/:id", middleware.RequireAuth, HomeBatteryController.GetConsumptionForSelectedDate)
		homeBatteryRoutes.POST("/status/selected-time/:id", middleware.RequireAuth, HomeBatteryController.GetStatusForSelectedTime)
		homeBatteryRoutes.POST("/status/selected-date/:id", middleware.RequireAuth, HomeBatteryController.GetStatusForSelectedDate)
	}
	uploadImageRoutes := r.Group("/api/upload")
	{
		imageUploadController := controllers.NewImageController()
		uploadImageRoutes.POST("/:real-estate-name", imageUploadController.Post)
		uploadImageRoutes.GET("/:file-name", imageUploadController.Get)
	}

	ambientSensor := r.Group("/api/ambient")
	{
		ambientSensorController := devicesController.NewAmbientSensorController(db, mqtt, cacheService)
		ambientSensor.GET("/:id", ambientSensorController.Get)
		ambientSensor.GET("/last-hour/:id", ambientSensorController.GetValueForHour)
		ambientSensor.POST("/selected-time/:id", ambientSensorController.GetValueForSelectedTime)
		ambientSensor.POST("/selected-date/:id", ambientSensorController.GetValuesForDate)
	}

	lampRoutes := r.Group("api/lamp")
	{
		lampController := devicesController.NewLampController(db, influxDb, cacheService)
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
		vehicleGateController := devicesController.NewVehicleGateController(db, influxDb, cacheService)
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
		vehicleGateRoutes.GET("/count/:id/:from/:to", vehicleGateController.GetEntriesOutcome)
	}

	permissionRoutes := r.Group("/api/permission")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		permissionController := controllers.NewPermissionController(db, cacheService)
		permissionRoutes.POST("", middleware.RequireAuth, permissionController.ReceiveGrantPermission)
		permissionRoutes.POST("/verify", permissionController.VerifyAccount)
		permissionRoutes.GET("/:id", middleware.RequireAuth, permissionController.GetPermissionForRealEstate)
		permissionRoutes.POST("/deny/:id", middleware.RequireAuth, permissionController.DeletePermit)
		permissionRoutes.GET("/get-real-estate/:id", middleware.RequireAuth, permissionController.GetPermitRealEstate)
		permissionRoutes.GET("/get-devices/:id/:userId", middleware.RequireAuth, permissionController.GetDeviceForRealEstate)
		permissionRoutes.GET("/get-permissions/:deviceId", middleware.RequireAuth, permissionController.GetPermissionsForDevice)
	}

	washingMachineRoutes := r.Group("/api/wm")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		washingMachineController := devicesController.NewWashingMachineController(db, mqtt, cacheService)
		washingMachineRoutes.GET("/:id", washingMachineController.Get)
		washingMachineRoutes.POST("/schedule", middleware.RequireAuth, washingMachineController.AddScheduledMode)
		washingMachineRoutes.GET("/schedule/:id", washingMachineController.GetScheduledModes)
		washingMachineRoutes.PUT("history", middleware.RequireAuth, washingMachineController.GetHistoryData)
	}

	electricityRoutes := r.Group("/api/consumption")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		ElectricityController := controllers.NewElectricityController(db, influxDb, &cacheService)
		electricityRoutes.POST("/selected-time", middleware.RequireAuth, ElectricityController.GetElectricityForSelectedTime)
		electricityRoutes.POST("/selected-date", middleware.RequireAuth, ElectricityController.GetElectricityForSelectedDate)
		electricityRoutes.POST("/ratio/selected-time", middleware.RequireAuth, ElectricityController.GetRatioForSelectedTime)
		electricityRoutes.POST("/ratio/selected-date", middleware.RequireAuth, ElectricityController.GetRatioForSelectedDate)
	}

	evChargerRoutes := r.Group("/api/ev")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		evChargerController := devicesController.NewEVChargerController(db, influxDb, cacheService)
		evChargerRoutes.GET("/:id", evChargerController.Get)
		evChargerRoutes.GET("/lastPercentage/:id", evChargerController.GetLastPercentage)
		evChargerRoutes.PUT("/actions", middleware.RequireAuth, evChargerController.GetHistoryActions)
	}
	SprinklerRoutes := r.Group("api/sprinkler")
	{
		middleware := middleware.NewMiddleware(db, cacheService)
		sprinklerController := devicesController.NewSprinklerController(db, influxDb, mqtt, cacheService)
		SprinklerRoutes.GET("/:id", sprinklerController.Get)
		SprinklerRoutes.GET("/", sprinklerController.GetAll)
		SprinklerRoutes.PUT("/:id/on", middleware.RequireAuth, sprinklerController.TurnOn)
		SprinklerRoutes.PUT("/:id/off", middleware.RequireAuth, sprinklerController.TurnOff)
		SprinklerRoutes.DELETE("/:id", sprinklerController.Delete)
		SprinklerRoutes.POST("/mode/:id", sprinklerController.AddSpecialMode)
		SprinklerRoutes.GET("/mode/:id", sprinklerController.GetSpecialModes)
		SprinklerRoutes.DELETE("/mode/:id", sprinklerController.DeleteSpecialMode)
		SprinklerRoutes.GET("/mode/one/:id", sprinklerController.GetSpecialMode)
		SprinklerRoutes.PUT("/history", middleware.RequireAuth, sprinklerController.GetHistoryData)
	}
}

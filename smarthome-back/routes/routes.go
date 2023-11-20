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

		// todo promeni middleware
		authController := controllers.NewAuthController(db)
		middleware := middleware.NewMiddleware(db)
		userRoutes.POST("/login", authController.Login)
		userRoutes.GET("/validate", middleware.RequireAuth, authController.Validate)
		userRoutes.POST("/logout", middleware.RequireAuth, authController.Logout)
		userRoutes.POST("/verificationMail", middleware.RequireAuth, authController.SendVerificationMail)
		userRoutes.POST("/activate", middleware.RequireAuth, authController.ActivateAccount)

		superadminController := controllers.NewSuperAdminController(db)
		userRoutes.POST("/reset-password", middleware.RequireAuth, superadminController.ResetPassword)
		userRoutes.POST("/add-admin", middleware.RequireAuth, superadminController.AddAdmin)
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
}

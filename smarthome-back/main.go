package main

import (
	_ "database/sql"
	_ "fmt"
	"smarthome-back/config"
	"smarthome-back/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	db := config.SetupDatabase()

	routes.SetupRoutes(r, db)

	// public := r.Group("/api")
	// public.POST("/register", controllers.Register)
	// public.POST("/login", controllers.Login)

	r.Run(":8081")
}

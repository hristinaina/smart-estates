package main

import (
	_ "database/sql"
	_ "fmt"
	"smarthome-back/config"
	"smarthome-back/routes"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.Use(config.SetupCORS())

	db := config.SetupDatabase()

	// Enable CORS for all routes
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r, db)

	gs := services.NewGenerateSuperAdmin(db)
	gs.GenerateSuperadmin()

	r.Run(":8081")
}

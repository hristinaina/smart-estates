package main

import (
	_ "database/sql"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/mqtt_client"
	"smarthome-back/routes"
)

func main() {
	r := gin.Default()
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

	mqttCLient := mqtt_client.NewMQTTClient()
	mqttCLient.StartListening()

	routes.SetupRoutes(r, db)
	r.Run(":8081")
}

package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
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
	resetDbOnStart(db)
	mqttClient := mqtt_client.NewMQTTClient(db)
	mqttClient.StartListening()

	routes.SetupRoutes(r, db)
	r.Run(":8081")
}

func resetDbOnStart(db *sql.DB) {
	query := "UPDATE device SET IsOnline = false"
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Failed to update devices status")
	}
}

package main

import (
	"database/sql"
	"fmt"
	"smarthome-back/config"
	"smarthome-back/mqtt_client"
	"smarthome-back/routes"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(config.SetupCORS())

	db := config.SetupDatabase()
	// session for aws
	// _, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("eu-central-1"),
	// 	Credentials: credentials.NewStaticCredentials("AKIAXTEDOKGSGESVDNWJ", "fXig4kJtKpMBK9q1NxGDpcVrm1xD+IqW1JeCOI7J", ""),
	// })

	//if err != nil {
	//	fmt.Println("Error while opening session on aws")
	//	panic(err)
	//}
	// Enable CORS for all routes
	// r.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}
	// 	c.Next()
	// })

	resetDbOnStart(db)
	mqttClient := mqtt_client.NewMQTTClient(db)
	mqttClient.StartListening()

	routes.SetupRoutes(r, db)

	gs := services.NewGenerateSuperAdmin(db)
	gs.GenerateSuperadmin()

	r.Run(":8081")
}

func resetDbOnStart(db *sql.DB) {
	query := "UPDATE device SET IsOnline = false"
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Failed to update devices status")
	}
}

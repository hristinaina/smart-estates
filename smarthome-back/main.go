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
	resetDbOnStart(db)

	mqttClient := mqtt_client.NewMQTTClient(db)
	if mqttClient == nil {
		fmt.Println("Failed to connect to mqtt broker")
	} else {
		mqttClient.StartListening()
		fmt.Println("Started listening to mqtt topics.")
	}
	routes.SetupRoutes(r, db, mqttClient)

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

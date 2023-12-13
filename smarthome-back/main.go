package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"smarthome-back/config"
	"smarthome-back/mqtt_client"
	"smarthome-back/routes"
	"smarthome-back/services"
)

func main() {
	r := gin.Default()
	r.Use(config.SetupCORS())

	db := config.SetupMySQL()
	// session for aws
	// _, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("eu-central-1"),
	// 	Credentials: credentials.NewStaticCredentials("AKIAXTEDOKGSGESVDNWJ", "fXig4kJtKpMBK9q1NxGDpcVrm1xD+IqW1JeCOI7J", ""),
	// })

	//if err != nil {
	//	fmt.Println("Error while opening session on aws")
	//	panic(err)
	//}

	influxDb, err := config.SetupInfluxDb()
	if err != nil {
		fmt.Println(err)
	}

	mqttClient := mqtt_client.NewMQTTClient(db, influxDb)
	if mqttClient == nil {
		fmt.Println("Failed to connect to mqtt broker")
	} else {
		mqttClient.StartListening()
		fmt.Println("Started listening to mqtt topics.")
	}

	routes.SetupRoutes(r, db, mqttClient, influxDb)
	gs := services.NewGenerateSuperAdmin(db)
	gs.GenerateSuperadmin()

	r.Run(":8081")
}

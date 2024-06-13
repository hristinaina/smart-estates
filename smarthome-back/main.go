package main

import (
	"fmt"
	"net/http"
	"smarthome-back/config"
	"smarthome-back/mqtt_client"
	"smarthome-back/routes"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adresa i port Redis servera
		Password: "",               // Lozinka, ako je postavljena
		DB:       0,                // Broj baze podataka, ako koristite više baza
	})

	influxDb, err := config.SetupInfluxDb()
	if err != nil {
		fmt.Println(err)
	}

	mqttClient := mqtt_client.NewMQTTClient(db, influxDb, redisClient)
	if mqttClient == nil {
		fmt.Println("Failed to connect to mqtt broker")
	} else {
		mqttClient.StartListening()
		fmt.Println("Started listening to mqtt topics.")
	}

	// web socket
	go func() {
		config.SetupWebSocketRoutes(db, influxDb)
		http.ListenAndServe(":8082", nil)
	}()

	routes.SetupRoutes(r, db, mqttClient, influxDb)
	gs := services.NewGenerateSuperAdmin(db)
	gs.GenerateSuperadmin()

	r.Run(":8081")

	defer db.Close()
}

package main

import (
	"fmt"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/config"
	"smarthome-back/mqtt_client"
	"smarthome-back/routes"
	"smarthome-back/services"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
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

	configCache := bigcache.Config{
		Shards:             1024,
		LifeWindow:         24 * time.Hour,
		CleanWindow:        10 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	cacheSetup, _ := bigcache.NewBigCache(configCache)
	cacheService := cache.NewCacheService(cacheSetup)

	influxDb, err := config.SetupInfluxDb()
	if err != nil {
		fmt.Println(err)
	}

	mqttClient := mqtt_client.NewMQTTClient(db, influxDb, cacheService)
	if mqttClient == nil {
		fmt.Println("Failed to connect to mqtt broker")
	} else {
		mqttClient.StartListening()
		fmt.Println("Started listening to mqtt topics.")
	}

	// web socket
	go func() {
		config.SetupWebSocketRoutes(db, influxDb, cacheService)
		http.ListenAndServe(":8082", nil)
	}()

	routes.SetupRoutes(r, db, mqttClient, influxDb, *cacheService)
	gs := services.NewGenerateSuperAdmin(db)
	gs.GenerateSuperadmin()

	r.Run(":8081")

	defer db.Close()
}

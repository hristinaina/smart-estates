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
		Shards:             1024,             // Adjust based on the number of concurrent accesses
		LifeWindow:         1 * time.Hour,    // Start with 1 hour and adjust based on your needs
		CleanWindow:        10 * time.Minute, // Regular clean-up to manage memory usage
		MaxEntriesInWindow: 1000 * 60,        // Adjust based on expected traffic and data access patterns
		MaxEntrySize:       1024,             // Adjust based on the size of your data entries
		Verbose:            false,            // Set to true for detailed logging during initial setup/debugging
		HardMaxCacheSize:   8192,             // Adjust based on available memory
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
	gs := services.NewGenerateSuperAdmin(db, *cacheService)
	gs.GenerateSuperadmin()

	r.Run(":8081")

	defer db.Close()
}

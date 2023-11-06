package main

import (
	"github.com/gin-gonic/gin"
	"smarthome-back/routes"
	"smarthome-back/config"
)


func main() {
	r := gin.Default()
	db := config.SetupDatabase()
	routes.SetupRoutes(r, db)
	r.Run(":8080")
}
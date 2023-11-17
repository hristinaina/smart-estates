package main

import (
	_ "database/sql"
	_ "fmt"
	"smarthome-back/config"
	"smarthome-back/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.Use(config.SetupCORS())

	db := config.SetupDatabase()

	routes.SetupRoutes(r, db)

	r.Run(":8081")
}

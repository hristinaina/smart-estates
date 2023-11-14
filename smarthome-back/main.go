package main

import (
	_ "database/sql"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/routes"
)

func main() {
	r := gin.Default()
	db := config.SetupDatabase()

	routes.SetupRoutes(r, db)
	r.Run(":8081")
}

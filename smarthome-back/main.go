package main

import (
	_ "database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/routes"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db := config.SetupDatabase()

	routes.SetupRoutes(r, db)
	err := r.Run(":8081")
	if err != nil {
		fmt.Println("Error happened: ", err)
		return
	}
}

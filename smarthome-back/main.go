package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/routes"
)

func main() {
	// Change the values of the connection string according to your MySQL setup
	// connected with scada database only for testing purposes
	database, err := sql.Open("mysql", "root:siit2020@tcp(localhost:3306)/scada")
	if err != nil {
		panic(err.Error())
	}
	defer database.Close()

	// Test the connection
	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected to the database!")

	r := gin.Default()
	db := config.SetupDatabase()

	routes.SetupRoutes(r, db, database)
	r.Run(":8081")
}

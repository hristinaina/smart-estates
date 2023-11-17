package main

import (
	_ "database/sql"
	_ "fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/routes"
	"smarthome-back/services"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db := config.SetupDatabase()

	r.GET("/send-email", func(c *gin.Context) {
		to := "kacorinav@gmail.com"

		if err := services.NewMailService(db).ApproveRealEstate(to); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Email sent successfully"})
	})

	routes.SetupRoutes(r, db)
	r.Run(":8081")
}

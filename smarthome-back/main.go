package main

import (
	_ "database/sql"
	"fmt"
	_ "fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
	"smarthome-back/config"
	"smarthome-back/routes"
)

var emailConfig = map[string]string{
	"host":     "smtp.gmail.com",
	"port":     "587",
	"username": "kacorinav@gmail.com",
	"password": "fhgm naqq bnkn dqxe",
}

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig["username"])
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailConfig["host"], 587, emailConfig["username"], emailConfig["password"])

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db := config.SetupDatabase()

	r.GET("/send-email", func(c *gin.Context) {
		to := "kacorinav@gmail.com"
		subject := "Test Email"
		body := "This is a test email sent from Gin and Go."

		if err := sendEmail(to, subject, body); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Email sent successfully"})
	})

	routes.SetupRoutes(r, db)
	r.Run(":8081")
}

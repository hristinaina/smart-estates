package main

import (
	_ "database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/config"
	"smarthome-back/routes"
)

func main() {
	// session for aws
	_, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials("AKIAXTEDOKGSGESVDNWJ", "fXig4kJtKpMBK9q1NxGDpcVrm1xD+IqW1JeCOI7J", ""),
	})

	if err != nil {
		fmt.Println("Error while opening session on aws")
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())
	db := config.SetupDatabase()
	routes.SetupRoutes(r, db)
	err = r.Run(":8081")
	if err != nil {
		fmt.Println("Error happened: ", err)
		return
	}
}

package config

import (
	"database/sql"
	"fmt"
)

func SetupDatabase() *sql.DB {

	database, err := sql.Open("mysql", "root:siit2020@tcp(localhost:3306)/smart_home")
	if err != nil {
		panic(err.Error())
	}
	//defer database.Close()

	// Test the connection
	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected to the database!")

	return database
}

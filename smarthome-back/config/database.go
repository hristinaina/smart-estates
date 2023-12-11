package config

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	_ "github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"smarthome-back/services"
	"time"
)

const (
	Org    = "Smart Home"
	Bucket = "bucket"
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

// SetupInfluxDb this function is the same as function main in influxdb module
// this is created just in case if we have problems with working in different modules
func SetupInfluxDb() (influxdb2.Client, error) {
	// connecting to database
	var service = services.NewConfigService()
	token, err := service.GetToken("config/config.json")
	if err != nil {
		fmt.Println("Error happened while connecting with InfluxDb!")
		return nil, err
	}

	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)

	return client, nil
}

func TestInfluxDb(client influxdb2.Client) {
	// writing data example
	// this is our organization
	writeAPI := client.WriteAPIBlocking(Org, Bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"Device Name": "Thermometer",
			"Position":    "Room 1",
		}
		fields := map[string]interface{}{
			"C": 22.5,
			"K": 295.65,
		}
		// measurement == table in relation db
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}

	// executing query
	// reading data example
	queryAPI := client.QueryAPI(Org)
	// we are printing data that came in the last 10 minutes
	query := `from(bucket: "bucket")
            |> range(start: -10m)
            |> filter(fn: (r) => r._measurement == "measurement1")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Printing data...")
	for results.Next() {
		fmt.Println("------------------------")
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

	// this is example of aggregate function (mean)
	// it only includes data that came in the last 10 minutes
	query = `from(bucket: "bucket")
              |> range(start: -10m)
              |> filter(fn: (r) => r._measurement == "measurement1")
              |> mean()`
	results, err = queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mean example:")
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}

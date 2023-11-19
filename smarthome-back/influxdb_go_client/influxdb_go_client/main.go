package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	_ "github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"time"
)

func main() {
	// connecting to mysql database
	token := ""
	fmt.Println("Variableee")
	fmt.Println(token)

	url := "http://127.0.0.1:8086"
	client := influxdb2.NewClient(url, token)

	// writing data example
	// this is our organization
	org := "Smart Home"
	bucket := "bucket"
	writeAPI := client.WriteAPIBlocking(org, bucket)
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
	queryAPI := client.QueryAPI(org)
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

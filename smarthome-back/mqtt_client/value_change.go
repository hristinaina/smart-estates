package mqtt_client

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"smarthome-back/enumerations"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"strconv"
	"strings"
	"time"
)

// TODO : create while loop that will do publishes

func (mc *MQTTClient) HandleValueChange(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("usaooooo")
	parts := strings.Split(msg.Topic(), "/")
	fmt.Println(msg.Topic())

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	device, err := mc.deviceRepository.Get(id)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("DEVICE LAMPPP: ", device.Id)
	percentage := string(msg.Payload())
	fmt.Println("Percentage: ", percentage)
	value, err := strconv.ParseFloat(percentage, 32)
	if err != nil {
		fmt.Println("Error converting string to float32:", err)
		return
	}

	// Convert float64 to float32
	val := float32(value)
	mc.CheckValue(val, device)
}

func (mc *MQTTClient) CheckValue(value float32, device models2.Device) {
	if value != device.LastValue {
		if device.Type == enumerations.Lamp {
			lamp, err := mc.lampRepository.Get(device.Id)
			if err != nil {
				fmt.Println("Error: ", err)
				// TODO: handle error
			} else {
				fmt.Println("Posting new lamp value to influxdb...")
				fmt.Println(lamp)
				mc.PostNewLampValue(lamp, value)
			}
		}
		// TODO: handle other type
	}
}

func (mc *MQTTClient) PostNewLampValue(lamp models.Lamp, percentage float32) {
	client := mc.influxDb
	writeAPI := client.WriteAPIBlocking("Smart Home", "bucket")
	tags := map[string]string{
		"Id":         strconv.Itoa(lamp.ConsumptionDevice.Device.Id),
		"DeviceName": lamp.ConsumptionDevice.Device.Name,
	}
	fields := map[string]interface{}{
		"Value": percentage,
	}
	// measurement == table in relation db
	point := write.NewPoint("measurement1", tags, fields, time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Posted to influx db...")
	// printing values from influx db (last 10 minutes)
	mc.GetLampsFromInfluxDb()

}

func (mc *MQTTClient) GetLampsFromInfluxDb() *api.QueryTableResult {
	client := mc.influxDb
	queryAPI := client.QueryAPI("Smart Home")
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

	return results
}

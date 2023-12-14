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

func (mc *MQTTClient) HandleValueChange(_ mqtt.Client, msg mqtt.Message) {
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
	percentage := string(msg.Payload())
	value, err := strconv.ParseFloat(percentage, 32)
	if err != nil {
		fmt.Println("Error converting string to float32:", err)
		return
	}

	// float64 to float32
	val := float32(value)
	mc.CheckValue(val, device)
}

func (mc *MQTTClient) CheckValue(value float32, device models2.Device) {
	if value != device.LastValue {
		if device.Type == enumerations.Lamp {
			lamp, err := mc.lampRepository.Get(device.Id)
			if err != nil {
				fmt.Println("Error happened: ", err)
				// TODO: handle error
			} else {
				if value != device.LastValue {
					fmt.Println("Posting new lamp value to influxdb...")
					fmt.Println(lamp)
					mc.PostNewLampValue(lamp, value)
					// after new value is added to influx, last value property needs to be updated
					_, err = mc.deviceRepository.UpdateLastValue(device.Id, value)
				}
			}
		}
		// TODO: handle other types
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
	point := write.NewPoint("lamps", tags, fields, time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Posted to influx db...")
	// printing values from influx
	//mc.GetLampsFromInfluxDb("2023-01-01T00:00:00Z", "2023-12-31T00:00:00Z")

}

// TODO: delete this later
func (mc *MQTTClient) GetLampsFromInfluxDb(from, to string) *api.QueryTableResult {
	client := mc.influxDb
	queryAPI := client.QueryAPI("Smart Home")
	// we are printing data that came in the last 10 minutes
	query := fmt.Sprintf(`from(bucket: "bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "lamps")`, from, to)
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

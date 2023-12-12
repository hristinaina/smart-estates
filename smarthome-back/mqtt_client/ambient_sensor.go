package mqtt_client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var switchOn = false

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func (mc *MQTTClient) ReceiveValue(client mqtt.Client, msg mqtt.Message) {
	if switchOn {
		payload := string(msg.Payload())

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		deviceId := int(data["id"].(float64))
		temperature := data["temperature"].(float64)
		humidity := data["humidity"].(float64)

		saveValueToInfluxDb(mc.influxDb, deviceId, temperature, humidity)

		fmt.Printf("Ambient Sensor, id=%v, temeprature: %v °C, humidity: %v %% \n", deviceId, temperature, humidity)
	}
}

func saveValueToInfluxDb(client influxdb2.Client, deviceId int, temperature, humidity float64) {
	fmt.Println("uslooo jeee")
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	point := influxdb2.NewPoint("measurement", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{
			"temperature": temperature,
			"humidity":    humidity,
		},
		time.Now()) // field

	// Write the point to InfluxDB
	writeAPI.WritePoint(point)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Ambient sensor influxdb")
}

func (mc *MQTTClient) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	switchOn = status == "true"
	fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
}

package mqtt_client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (mc *MQTTClient) HandleActionChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	payload := string(msg.Payload())

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	mode := data["Mode"]
	switchAC := data["Switch"].(bool)
	temp := data["Temp"].(float64)
	previous := data["Previous"]
	user := data["UserEmail"]
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(deviceId, mode, temp, previous, user, switchAC)

	// todo sacuvaj u bazi
	// switchOn = status == "true"
	// fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
}

// func saveACToInfluxDb(client influxdb2.Client, deviceId int, mode, previous, user string, temp float64) {
// 	Org := "Smart Home"
// 	Bucket := "bucket"
// 	writeAPI := client.WriteAPI(Org, Bucket)

// 	point := influxdb2.NewPoint("air_conditioner", // table
// 		map[string]string{"device_id": strconv.Itoa(deviceId), "user_id": user}, // tag
// 		map[string]interface{}{"temperature": temperature, "humidity": humidity},
// 		time.Now()) // field

// 	// Write the point to InfluxDB
// 	writeAPI.WritePoint(point)

// 	// Close the write API to flush the buffer and release resources
// 	writeAPI.Flush()
// 	fmt.Println("Ambient sensor influxdb")
// }

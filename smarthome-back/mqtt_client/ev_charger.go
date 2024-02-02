package mqtt_client

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	"strconv"
	"strings"
	"time"
)

func (mc *MQTTClient) HandleInputPercentageChange(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the JSON string into the struct
	var data dtos.ElectricVehicleDTO
	err = json.Unmarshal([]byte(msg.Payload()), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	fmt.Println(data.Email)
	saveActionToInflux(mc.influxDb, deviceId, data.Email, "percentageChange", -1, data.CurrentCapacity)
}

func (mc *MQTTClient) HandleAutoActionsForCharger(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	fmt.Println(deviceId)
	if err != nil {
		fmt.Println(err)
	}

	var data dtos.ElectricVehicleDTO
	err = json.Unmarshal([]byte(msg.Payload()), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	saveActionToInflux(mc.influxDb, deviceId, "auto", data.Action, data.PlugId, data.CurrentCapacity)
}

// actions can be start, end and percentageChange
func saveActionToInflux(client influxdb2.Client, deviceId int, user string, action string, plugId int, percentage float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("ev_charger", //table
		map[string]string{"device_id": strconv.Itoa(deviceId), "user_id": user, "plug_id": strconv.Itoa(plugId), "action": action}, //tag
		map[string]interface{}{"value": percentage}, //field
		time.Now())

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Printf("Savin to influxdb. Electrical charger: id=%d, action %s, plugId %d, percentage %f \n", deviceId, action, plugId, percentage)

}

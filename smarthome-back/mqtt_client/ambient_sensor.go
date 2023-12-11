package mqtt_client

import (
	"encoding/json"
	"fmt"
	models "smarthome-back/models/devices"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func (mc *MQTTClient) ReceiveValue(client mqtt.Client, msg mqtt.Message) {
	// Ovde manipulišite podacima koji su stigli u poruci
	payload := string(msg.Payload())
	fmt.Println("Received message:", payload)

	// Primer kako da dekodirate JSON payload
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Sada možete koristiti vrednosti iz mape "data" kako vam odgovara
	deviceId := data["id"].(float64)
	temperature := data["temperature"].(float64)
	humidity := data["humidity"].(float64)

	// Radite šta god želite sa ovim vrednostima...
	fmt.Printf("Device is online, id=%v, temeprature: %v °C, humidity: %v %% \n", deviceId, temperature, humidity)
}

// CheckDeviceStatus function that checks if there is a device that has disconnected
// func (mc *MQTTClient) CheckDeviceStatus() {
// 	offlineTimeout := 30 * time.Second
// 	devices := mc.deviceRepository.GetAll()
// 	for _, device := range devices {
// 		if device.IsOnline && time.Since(device.StatusTimeStamp.Time) > offlineTimeout {
// 			fmt.Printf("Device with id=%d is offline.\n", device.Id)
// 			device.IsOnline = false
// 			device.StatusTimeStamp = mysql.NullTime{
// 				Time:  time.Now(),
// 				Valid: true,
// 			}
// 			mc.deviceRepository.Update(device)
// 			err := mc.Publish(TopicStatusChanged+strconv.Itoa(device.Id), "offline")
// 			saveToInfluxDb(mc.influxDb, device)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}
// }

func saveValueToInfluxDb(client influxdb2.Client, device models.Device) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	p := influxdb2.NewPoint("device_status", //table
		map[string]string{"device_id": strconv.Itoa(device.Id)}, //tag
		map[string]interface{}{"status": func() int {
			if device.IsOnline {
				return 1
			} else {
				return 0
			}
		}()}, //field
		device.StatusTimeStamp.Time)

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved status change to influxdb")
}

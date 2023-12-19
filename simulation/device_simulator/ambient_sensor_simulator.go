package device_simulator

import (
	"encoding/json"
	"fmt"
	"simulation/config"
	"simulation/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic = "ambientSensor/switch/"
)

type AmbientSensorSimulator struct {
	switchOn bool
	client   mqtt.Client
	device   models.AmbientSensor
}

func NewAmbientSensorSimulator(client mqtt.Client, device models.Device) *AmbientSensorSimulator {
	as, err := config.GetAmbientSensor(device.ID)
	if err != nil {
		return nil
	}
	return &AmbientSensorSimulator{
		client:   client,
		device:   as,
		switchOn: false,
	}
}

func (as *AmbientSensorSimulator) ConnectAmbientSensor() {
	go SendHeartBeat(as.client, as.device.Device.ID, as.device.Device.Name)
	go as.GenerateAmbientSensorData()
	// config.SubscribeToTopic(as.client, topicSwitch+strconv.Itoa(as.device.ID), as.HandleSwitchChange)
}

// za back
// GenerateAmbientSensorData Simulate sending periodic AmbientSensor data
func (as *AmbientSensorSimulator) GenerateAmbientSensorData() {
	var indoorTemperature, indoorHumidity float64
	slope := 0.5
	intercept := 15.0

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			openMeteoResponse, err := config.GetTemp()
			if err != nil {
				fmt.Printf("Error: %v \n", err.Error())
			}

			indoorTemperature = slope*openMeteoResponse.Current.Temperature2m + intercept
			indoorHumidity = openMeteoResponse.Current.RelativeHumidity2m / 2

			data := map[string]interface{}{
				"id":          as.device.Device.ID,
				"temperature": indoorTemperature,
				"humidity":    indoorHumidity,
			}
			jsonString, err := json.Marshal(data)
			if err != nil {
				fmt.Println("greska")
			}
			config.PublishToTopic(as.client, "device/ambient/sensor", string(jsonString)) // todo eventualno promeni topic ako bude potrebno
			fmt.Printf("AmbientSensor name=%s, id=%d, temeprature: %v °C, humidity: %v %% \n", as.device.Device.Name, as.device.Device.ID, indoorTemperature, indoorHumidity)
		}
	}
}

// za front
// func (as *AmbientSensorSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
// 	parts := strings.Split(msg.Topic(), "/")
// 	deviceId, err := strconv.Atoi(parts[len(parts)-1])
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	status := string(msg.Payload())
// 	as.switchOn = status == "true"
// 	fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
// }

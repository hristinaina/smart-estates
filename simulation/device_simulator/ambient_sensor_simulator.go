package device_simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
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
		switchOn: true,
	}
}

func (as *AmbientSensorSimulator) ConnectAmbientSensor() {
	go SendHeartBeat(as.client, as.device.Device.ID, as.device.Device.Name)
	go as.GenerateAmbientSensorData()
	go as.SendConsumption()
	//config.SubscribeToTopic(as.client, topicSwitch+strconv.Itoa(as.device.ID), as.HandleSwitchChange)
}

func (as *AmbientSensorSimulator) SendConsumption() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			scalingFactor := 1.0
			if as.switchOn {
				scalingFactor = 0.8 + rand.Float64()*0.2 // get a number between 0.8 and 1.0
			} else {
				scalingFactor = 0.15 + rand.Float64()*0.2 // get a number between 0.15 and 0.35
			}
			consumed := as.device.PowerConsumption * scalingFactor / 60 / 2 // divide by 60 and 2 to get consumption for previous 30s
			err := config.PublishToTopic(as.client, config.TopicConsumption+strconv.Itoa(as.device.Device.ID), strconv.FormatFloat(consumed,
				'f', -1, 64))
			if err != nil {
				fmt.Printf("Error publishing message with the device: %s \n", as.device.Device.Name)
			} else {
				fmt.Printf("%s: Ambient Sensor with id=%d, Name=%s, consumed=%fkWh for previous 30s\n", time.Now().Format("15:04:05"),
					as.device.Device.ID, as.device.Device.Name, consumed)
			}
		}
	}
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

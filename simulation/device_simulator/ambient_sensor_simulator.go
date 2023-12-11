package device_simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic = "ambientSensor/switch/"
)

type AmbientSensorSimulator struct {
	switchOn bool
	client   mqtt.Client
	device   models.Device
}

func NewAmbientSensorSimulator(client mqtt.Client, device models.Device) *AmbientSensorSimulator {
	//todo da se proslijedi samo deviceId (umjesto device) i posalje upit ka beku za dobavljane svih podataka za AmbientSensoru
	// (jer device ima samo opste podatke)
	return &AmbientSensorSimulator{
		client:   client,
		device:   device,
		switchOn: false,
	}
}

func (as *AmbientSensorSimulator) ConnectAmbientSensor() {
	go SendHeartBeat(as.client, as.device)
	go as.GenerateAmbientSensorData()
	config.SubscribeToTopic(as.client, topicSwitch+strconv.Itoa(as.device.ID), as.HandleSwitchChange)
}

// za back
// GenerateAmbientSensorData Simulate sending periodic AmbientSensor data
func (as *AmbientSensorSimulator) GenerateAmbientSensorData() {
	temperature := 22
	humidity := 35

	for {
		// if as.switchOn {
		temperature = temperature + rand.Intn(3) - 1
		humidity = humidity + rand.Intn(3) - 1
		if humidity < 0 {
			humidity = 0
		}
		if humidity > 100 {
			humidity = 100
		}
		data := map[string]interface{}{
			"id":          as.device.ID,
			"temperature": temperature,
			"humidity":    humidity,
		}
		jsonString, err := json.Marshal(data)
		if err != nil {
			fmt.Println("greska")
		}
		config.PublishToTopic(as.client, "device/ambient/sensor", string(jsonString)) // todo eventualno promeni topic ako bude potrebno
		fmt.Printf("AmbientSensor name=%s, id=%d, temeprature: %v °C, humidity: %v %% \n", as.device.Name, as.device.ID, temperature, humidity)
		time.Sleep(5 * time.Second)
		// }
	}
}

// za front
func (as *AmbientSensorSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	as.switchOn = status == "true"
	fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
}

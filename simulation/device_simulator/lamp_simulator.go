package device_simulator

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"simulation/config"
	"simulation/models"
	"time"
)

func ConnectLamp(client mqtt.Client, device models.Device) {
	go config.SendHeartBeat(client, device)
	go GenerateLampData(client)
}

// GenerateLampData Simulate sending periodic Lamp data
func GenerateLampData(client mqtt.Client) {
	for {
		config.SendMessage(client, config.TopicPayload, "some simulated data")
		time.Sleep(5 * time.Second)
	}
}

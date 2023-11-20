package device_simulator

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"simulation/config"
	"simulation/models"
	"strconv"
	"time"
)

func ConnectLamp(client mqtt.Client, device models.Device) {
	go config.SendHeartBeat(client, device)
	go GenerateLampData(client, device)
}

// GenerateLampData Simulate sending periodic Lamp data
func GenerateLampData(client mqtt.Client, device models.Device) {
	for {
		config.SendMessage(client, config.TopicPayload+strconv.Itoa(device.ID), "some simulated data")
		time.Sleep(5 * time.Second)
	}
}

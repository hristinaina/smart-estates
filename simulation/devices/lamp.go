package devices

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func ConnectLamp(client mqtt.Client) {
	go SendHeartBeat(client)
	go GenerateLampData(client)
}

// GenerateLampData Simulate sending periodic Lamp data
func GenerateLampData(client mqtt.Client) {
	for {
		SendMessage(client, topicPayload, "some simulated data")
		time.Sleep(5 * time.Second)
	}
}

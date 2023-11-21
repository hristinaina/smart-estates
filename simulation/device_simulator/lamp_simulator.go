package device_simulator

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math"
	"simulation/config"
	"simulation/models"
	"strconv"
	"time"
)

func ConnectLamp(client mqtt.Client, device models.Device) {
	go SendHeartBeat(client, device)
	go GenerateLampData(client, device)
}

// GenerateLampData Simulate sending periodic Lamp data
func GenerateLampData(client mqtt.Client, device models.Device) {
	for {
		// Get the Unix timestamp from the current time
		unixTimestamp := float64(time.Now().Unix())
		sineValue := math.Sin(unixTimestamp)
		percentage := math.Abs(math.Round(sineValue * 100))
		config.PublishToTopic(client, config.TopicPayload+strconv.Itoa(device.ID), strconv.FormatFloat(percentage,
			'f', -1, 64))
		fmt.Printf("Lamp name=%s, id=%d, generated data: %f\n", device.Name, device.ID, percentage)
		time.Sleep(5 * time.Second)
	}
}

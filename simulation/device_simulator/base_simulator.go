package device_simulator

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"
)

// SendHeartBeat Periodically send online status
func SendHeartBeat(client mqtt.Client, device models.Device) {
	for {
		err := config.PublishToTopic(client, config.TopicOnline+strconv.Itoa(device.ID), "online")
		if err != nil {
			fmt.Printf("Error publishing message with the device: %s \n", device.Name)
		} else {
			fmt.Printf("%s: Device sent a heartbeat, id=%d, Name=%s \n", time.Now().Format("15:04:05"),
				device.ID, device.Name)
		}
		time.Sleep(10 * time.Second)
	}
}

func HandleNewDevice(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Printf("Error getting device with id: %d \n", deviceId)
		return
	}
	device, err := config.Get(deviceId)
	if err != nil {
		fmt.Println(err)
		return
	}
	StartSimulation(client, device)
}

func StartSimulation(client mqtt.Client, d models.Device) {
	switch d.Type {
	case 0:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		lamp := NewLampSimulator(client, d)
		go lamp.ConnectLamp()
	case 1:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		lamp := NewLampSimulator(client, d)
		go lamp.ConnectLamp()
	default:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		lamp := NewLampSimulator(client, d)
		go lamp.ConnectLamp()
		//todo change this and add separate logic for each device type
	}
}

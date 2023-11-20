package config

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"simulation/models"
	"strconv"
	"time"
)

const (
	TopicOnline  = "device/online/" //device/online/{deviceId}
	TopicPayload = "device/data/"   //device/data/{deviceId}
)

func CreateConnection() mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("ws://broker.emqx.io:8083/mqtt")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return client
}

// SendHeartBeat Periodically send online status
func SendHeartBeat(client mqtt.Client, device models.Device) {
	for {
		err := SendMessage(client, TopicOnline+strconv.Itoa(device.ID), "online")
		if err != nil {
			fmt.Printf("Error publishing message with the device: %s \n", device.Name)
		} else {
			fmt.Printf("%s: Device sent a heartbeat, id=%d, Name=%s \n", time.Now().Format("15:04:05"),
				device.ID, device.Name)
		}
		time.Sleep(10 * time.Second)
	}
}

func SendMessage(client mqtt.Client, topic, message string) error {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

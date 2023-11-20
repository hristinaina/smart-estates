package config

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

const (
	TopicOnline  = "device/online"
	TopicPayload = "device/data"
)

func CreateConnection() mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return client
}

// SendHeartBeat Periodically send online status
func SendHeartBeat(client mqtt.Client) {
	for {
		SendMessage(client, TopicOnline, "online")
		time.Sleep(10 * time.Second)
	}
}

func SendMessage(client mqtt.Client, topic, message string) {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
}

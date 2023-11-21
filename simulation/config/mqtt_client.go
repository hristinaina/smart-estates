package config

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

const (
	TopicOnline    = "device/online/" //device/online/{deviceId}
	TopicPayload   = "device/data/"   //device/data/{deviceId}
	TopicNewDevice = "device/new/"    //device/data/{deviceId}
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

// PublishToTopic method to publish data to topic
func PublishToTopic(client mqtt.Client, topic, message string) error {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

// SubscribeToTopic method to subscribe to topic
func SubscribeToTopic(client mqtt.Client, topic string, handler mqtt.MessageHandler) {
	token := client.Subscribe(topic, 1, handler)
	token.Wait()

	// Check if the subscription was successful
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

package config

import (
	"fmt"
	"os"
	"time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TopicOnline    = "device/online/" //device/online/{deviceId}
	TopicPayload   = "device/data/"   //device/data/{deviceId}
	TopicNewDevice = "device/new/"    //device/data/{deviceId}
)

func CreateConnection() mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("ws://localhost:9001/mqtt")
	opts.SetClientID("go-simulator-nvt-2023")
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)

		// Attempt to reconnect
		for {
			fmt.Println("Attempting to reconnect...")
			token := client.Connect()
			if token.Wait() && token.Error() == nil {
				fmt.Println("Reconnected successfully!")
				break
			}

			// Wait before attempting again
			time.Sleep(5 * time.Second)
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return client
}

// PublishToTopic method to publish data to topic
func PublishToTopic(client mqtt.Client, topic, message string) error {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

// SubscribeToTopic method to subscribe to topic
func SubscribeToTopic(client mqtt.Client, topic string, handler mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, handler)
	token.Wait()

	// Check if the subscription was successful
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

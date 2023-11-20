package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type MQTTClient struct {
	client mqtt.Client
	//todo add influxdb
}

/*
	   Topics:  	device/online/+ to subscribe to all messages
					device/online/{deviceId} to subscribe to only with that id
*/
const (
	topicOnline  = "device/online/"
	topicPayload = "device/data/"
)

func NewMQTTClient() *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("ws://broker.emqx.io:8083/mqtt")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return &MQTTClient{
		client: client,
		//todo add influxdb
	}
}

func (mc *MQTTClient) StartListening() {
	mc.subscribeToTopic(topicOnline+"+", HandleHeartBeat)
	//todo subscribe to all topics here and create your callback function otherplace (in other file)

	// Periodically check if the device is still online
	go func() {
		for {
			CheckDeviceStatus(mc.client)
			time.Sleep(15 * time.Second)
		}
	}()
}

func (mc *MQTTClient) subscribeToTopic(topic string, handler mqtt.MessageHandler) {
	token := mc.client.Subscribe(topic, 1, handler)
	token.Wait()

	// Check if the subscription was successful
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
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

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
device/online/+ to subscribe to all messages
device/online/{deviceId} to subscribe to only ones with that id
*/
const (
	TopicOnline    = "device/online/"
	TopicPayload   = "device/data/"
	TopicNewDevice = "device/new/"
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
	mc.SubscribeToTopic(TopicOnline+"+", HandleHeartBeat)
	//todo subscribe here to other topics. Create your callback functions in other file

	// Periodically check if the device is still online
	go func() {
		for {
			CheckDeviceStatus(mc.client)
			time.Sleep(15 * time.Second)
		}
	}()
}

func (mc *MQTTClient) SubscribeToTopic(topic string, handler mqtt.MessageHandler) {
	token := mc.client.Subscribe(topic, 1, handler)
	token.Wait()

	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// Publish this function is the same as PublishToTopic, the only difference is that it belongs to MQTTClient class
func (mc *MQTTClient) Publish(topic, message string) error {
	token := mc.client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

// PublishToTopic this function is the same as Publish, the only difference is that it doesn't belong to MQTTClient class
func PublishToTopic(client mqtt.Client, topic, message string) error {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

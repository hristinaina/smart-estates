package mqtt_client

import (
	"database/sql"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

/*
device/online/+ to subscribe to all messages
device/online/{deviceId} to subscribe to only ones with that id
*/
const (
	TopicOnline        = "device/online/" //from simulation to back
	TopicStatusChanged = "device/status/" //from back to front (because front doesn't have to know about everything from simulation)
	TopicPayload       = "device/data/"
	TopicNewDevice     = "device/new/"
)

type MQTTClient struct {
	client mqtt.Client
	db     *sql.DB
	//todo add influxdb
}

func NewMQTTClient(db *sql.DB) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("ws://broker.emqx.io:8083/mqtt")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return &MQTTClient{
		client: client,
		db:     db,
		//todo add influxdb
	}
}

func (mc *MQTTClient) StartListening() {
	mc.SubscribeToTopic(TopicOnline+"+", mc.HandleHeartBeat)
	//todo subscribe here to other topics. Create your callback functions in other file

	// Periodically check if the device is still online
	go func() {
		for {
			mc.CheckDeviceStatus()
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

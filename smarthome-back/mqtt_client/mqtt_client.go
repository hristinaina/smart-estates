package mqtt_client

import (
	"database/sql"
	"fmt"
	"os"
	"smarthome-back/repositories"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	TopicAmbientSensor = "device/ambient/sensor"
)

type MQTTClient struct {
	client           mqtt.Client
	deviceRepository repositories.DeviceRepository
	influxDb         influxdb2.Client
}

func NewMQTTClient(db *sql.DB, influxDb influxdb2.Client) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("ws://localhost:9001/mqtt")
	opts.SetClientID("go-server-nvt-2023")
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
		return nil
	}
	return &MQTTClient{
		client:           client,
		deviceRepository: repositories.NewDeviceRepository(db),
		influxDb:         influxDb,
	}
}

func (mc *MQTTClient) StartListening() {
	mc.SubscribeToTopic(TopicOnline+"+", mc.HandleHeartBeat)
	mc.SubscribeToTopic("lamp/switch/"+"+", mc.HandleSwitchChange)
	mc.SubscribeToTopic(TopicAmbientSensor, mc.ReceiveValue)
	//todo subscribe here to other topics. Create your callback functions in other file

	// Periodically check if the device is still online
	go func() {
		for {
			fmt.Println("checking device status...")
			mc.CheckDeviceStatus()
			time.Sleep(15 * time.Second)
		}
	}()
}

func (mc *MQTTClient) SubscribeToTopic(topic string, handler mqtt.MessageHandler) {
	token := mc.client.Subscribe(topic, 0, handler)
	token.Wait()

	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("Subscribed to topic = " + topic)
}

func (mc *MQTTClient) Publish(topic string, message string) error {
	token := mc.client.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Error publishing message:", token.Error())
	}
	return token.Error()
}

func (mc *MQTTClient) GetInflux() influxdb2.Client {
	return mc.influxDb
}

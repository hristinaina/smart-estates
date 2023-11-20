package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type MQTTClient struct {
	client     mqtt.Client
	heartBeats map[int]time.Time
	//todo add db influxdb
}

const (
	topicOnline  = "device/online"
	topicPayload = "device/data"
)

func NewMQTTClient() *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return &MQTTClient{
		client:     client,
		heartBeats: make(map[int]time.Time),
		//todo add db
	}
}

func (mc *MQTTClient) StartListening() {
	mc.subscribeToTopic(topicOnline, mc.handleHeartBeat)
	//todo subscribe to all topics here and create your callback function otherplace (in other file)

	// Periodically check if the device is still online
	go func() {
		for {
			mc.checkDeviceStatus()
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

func (mc *MQTTClient) handleHeartBeat(client mqtt.Client, msg mqtt.Message) {
	// Update heartbeat time when online message is received
	deviceId := int(msg.Payload()[0])
	mc.heartBeats[deviceId] = time.Now()
	fmt.Printf("Device is online, id=%d\n", deviceId)
}

func (mc *MQTTClient) checkDeviceStatus() {
	offlineTimeout := 30 * time.Second

	for deviceID, timestamp := range mc.heartBeats {
		if time.Since(timestamp) > offlineTimeout {
			fmt.Printf("Device with id=%d is offline.\n", deviceID)
			// todo posalji frontu da je uredjaj offline
			// mozemo ga izbaciti iz liste da ne bi stalno slao?
		}
	}
}

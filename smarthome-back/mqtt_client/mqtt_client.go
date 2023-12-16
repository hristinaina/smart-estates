package mqtt_client

import (
	"database/sql"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"os"
	"smarthome-back/repositories"
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
	TopicSPSwitch      = "sp/switch/"
	TopicSPData        = "sp/data/"
)

type MQTTClient struct {
	client               mqtt.Client
	deviceRepository     repositories.DeviceRepository
	solarPanelRepository repositories.SolarPanelRepository
	realEstateRepository repositories.RealEstateRepository
	influxDb             influxdb2.Client
}

func NewMQTTClient(db *sql.DB, influxDb influxdb2.Client) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("ws://localhost:9001/mqtt")
	opts.SetClientID("go-server-nvt-2023")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return nil
	}
	return &MQTTClient{
		client:               client,
		deviceRepository:     repositories.NewDeviceRepository(db),
		solarPanelRepository: repositories.NewSolarPanelRepository(db),
		realEstateRepository: *repositories.NewRealEstateRepository(db),
		influxDb:             influxDb,
	}
}

func (mc *MQTTClient) StartListening() {
	mc.SubscribeToTopic(TopicOnline+"+", mc.HandleHeartBeat)
	mc.SubscribeToTopic("lamp/switch/"+"+", mc.HandleSwitchChange)
	mc.SubscribeToTopic(TopicAmbientSensor, mc.ReceiveValue)
	mc.SubscribeToTopic(TopicSPSwitch+"+", mc.HandleSPSwitch)
	mc.SubscribeToTopic(TopicSPData+"+", mc.HandleSPData)
	//todo subscribe here to other topics. Create your callback functions in other file

	mc.StartConsumptionThread()
	mc.StartDeviceStatusThread()
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

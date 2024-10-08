package mqtt_client

import (
	"database/sql"
	"fmt"
	"os"
	"smarthome-back/cache"
	"smarthome-back/repositories"
	repositories2 "smarthome-back/repositories/devices"

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
	TopicSPSwitch      = "sp/switch/"
	TopicSPData        = "sp/data/"
	TopicACAction      = "ac/action/"
	TopicACSwitch      = "ac/switch/"
	TopicWMSwitch      = "wm/switch/"
	TopicConsumption   = "device/consumption/"
	TopicApproached    = "device/approached/"
	TopicVGOpenClose   = "vg/open/"
	TopicEVCStart      = "ev/start/"
	TopicEVCEnd        = "ev/end/"
	TopicEVCPercentage = "ev/percentage/"
	TurnSprinklerON    = "sprinkler/on/"
	TurnSprinklerOFF   = "sprinkler/off/"
)

type MQTTClient struct {
	client                mqtt.Client
	deviceRepository      repositories2.DeviceRepository
	solarPanelRepository  repositories2.SolarPanelRepository
	lampRepository        repositories2.LampRepository
	influxDb              influxdb2.Client
	realEstateRepository  repositories.RealEstateRepository
	homeBatteryRepository repositories2.HomeBatteryRepository
	vehicleGateRepository repositories2.VehicleGateRepository
	evChargerRepository   repositories2.EVChargerRepository
	sprinkleRepository    repositories2.SprinklerRepository
	cacheService          cache.CacheService
}

func NewMQTTClient(db *sql.DB, influxDb influxdb2.Client, cacheService *cache.CacheService) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker("ws://localhost:9001/mqtt")
	opts.SetClientID("go-server-nvt-2023")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return nil
	}
	return &MQTTClient{
		client:                client,
		deviceRepository:      repositories2.NewDeviceRepository(db, influxDb, cacheService),
		solarPanelRepository:  repositories2.NewSolarPanelRepository(db, *cacheService),
		lampRepository:        repositories2.NewLampRepository(db, influxDb, *cacheService),
		homeBatteryRepository: repositories2.NewHomeBatteryRepository(db, *cacheService),
		realEstateRepository:  *repositories.NewRealEstateRepository(db, cacheService),
		vehicleGateRepository: repositories2.NewVehicleGateRepository(db, influxDb, *cacheService),
		evChargerRepository:   repositories2.NewEVChargerRepository(db, *cacheService),
		sprinkleRepository:    repositories2.NewSprinklerRepository(db, influxDb, *cacheService),
		influxDb:              influxDb,
	}
}

func (mc *MQTTClient) StartListening() {
	mc.SubscribeToTopic(TopicOnline+"+", mc.HandleHeartBeat)
	mc.SubscribeToTopic(TopicAmbientSensor, mc.ReceiveValue)
	mc.SubscribeToTopic(TopicSPSwitch+"+", mc.HandleSPSwitch)
	mc.SubscribeToTopic(TopicSPData+"+", mc.HandleSPData)
	mc.SubscribeToTopic(TopicACSwitch+"+", mc.HandleActionChange)
	mc.SubscribeToTopic(TopicWMSwitch+"+", mc.HandleWMAction)
	mc.SubscribeToTopic(TopicPayload+"+", mc.HandleValueChange)
	mc.SubscribeToTopic(TopicConsumption+"+", mc.HandleHBData)
	mc.SubscribeToTopic(TopicApproached+"+", mc.HandleVehicleApproached)
	mc.SubscribeToTopic(TopicEVCStart+"+", mc.HandleAutoActionsForCharger)
	mc.SubscribeToTopic(TopicEVCEnd+"+", mc.HandleAutoActionsForCharger)
	mc.SubscribeToTopic(TopicEVCPercentage+"+", mc.HandleInputPercentageChange)
	mc.SubscribeToTopic(TurnSprinklerON+"+", mc.HandleSprinklerMessage)
	mc.SubscribeToTopic(TurnSprinklerOFF+"+", mc.HandleSprinklerOffMessage)
	//todo subscribe here to other topics. Create your callback functions in other file

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

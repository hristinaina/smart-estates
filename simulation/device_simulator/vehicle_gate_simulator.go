package device_simulator

import (
	// "math"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicApproached = "device/approached/"
	TopicVGOpenClose = "vg/open/"
)

type VehicleGateSimulator struct {
	client mqtt.Client
	device models.Device
}

func NewVehicleGateSimulator(client mqtt.Client, device models.Device) *VehicleGateSimulator{
	return &VehicleGateSimulator {
		client: client,
		device: device,
	}
}

func (sim *VehicleGateSimulator) ConnectVehicleGate() {
	go SendHeartBeat(sim.client, sim.device.ID, sim.device.Name)
	go sim.GenerateVehicleData()
	go config.SubscribeToTopic(sim.client, TopicVGOpenClose+strconv.Itoa(sim.device.ID), sim.HandleLeaving)
}

func (sim *VehicleGateSimulator) GenerateVehicleData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Float64()
			if randomNumber <= 0.18 {
				sim.HandleCarApproached()
			}
		}
	}
}

func (sim *VehicleGateSimulator) HandleCarApproached() {
	licensePlate := "NS-123-45"
	config.PublishToTopic(sim.client, config.TopicApproached+strconv.Itoa(sim.device.ID), licensePlate+"+enter")
	fmt.Println("Published to topic approached!")
}

func (sim *VehicleGateSimulator) HandleLeaving(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Someone is leaving\n")
	payload := string(msg.Payload())
	payloadTokens := strings.Split(payload, "+")
	fmt.Printf("payload: %s\n", payload)
	if (payloadTokens[0] == "open") {
		licensePlate := payloadTokens[1]
		action := payloadTokens[2]
		if action == "enter" {
			fmt.Printf("Simulation %s is leaving...\n", licensePlate)
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Float64()
			sec := int(randomNumber * 30)
			fmt.Printf("Leaving in %d seconds\n", sec)
			timerChan := time.After(time.Duration(sec) * time.Second)
			select {
			case <- timerChan:
				fmt.Printf("Left %s\n", licensePlate)
				config.PublishToTopic(sim.client, config.TopicApproached+strconv.Itoa(sim.device.ID), licensePlate+"+exit")
			}
			
		}
	}

}
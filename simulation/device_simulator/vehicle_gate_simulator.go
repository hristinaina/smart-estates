package device_simulator

import (
	// "math"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
	// "strings"
	"time"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicApproached = "device/approached/"
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
}

func (sim *VehicleGateSimulator) GenerateVehicleData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Float64()
			if randomNumber <= 0.3 {
				sim.HandleCarApproached()
			}
		}
	}
}

func (sim *VehicleGateSimulator) HandleCarApproached() {
	licensePlate := "NS-123-45"
	config.PublishToTopic(sim.client, config.TopicApproached+strconv.Itoa(sim.device.ID), licensePlate)
	fmt.Println("Published to topic approached!")
}
package device_simulator

import (
	"simulation/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BatterySimulator struct {
	client mqtt.Client
	device models.Device
}

func NewBatterySimulator(client mqtt.Client, device models.Device) *BatterySimulator {
	return &BatterySimulator{
		client: client,
		device: device,
	}
}

func (ls *BatterySimulator) ConnectBattery() {
	go SendHeartBeat(ls.client, ls.device.ID, ls.device.Name)
}

package device_simulator

import (
	"fmt"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SendHeartBeat Periodically send online status
func SendHeartBeat(client mqtt.Client, id int, name string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := config.PublishToTopic(client, config.TopicOnline+strconv.Itoa(id), "online")
			if err != nil {
				fmt.Printf("Error publishing message with the device: %s \n", name)
			} else {
				fmt.Printf("%s: Device sent a heartbeat, id=%d, Name=%s \n", time.Now().Format("15:04:05"),
					id, name)
			}
		}
	}
}

func HandleNewDevice(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Printf("Error getting device with id: %d \n", deviceId)
		return
	}
	device, err := config.Get(deviceId)
	if err != nil {
		fmt.Println(err)
		return
	}
	StartSimulation(client, device)
}

func StartSimulation(client mqtt.Client, d models.Device) {
	// fmt.Printf("DEviceeeee: %s\n", d.Name)
	// fmt.Printf("Typeee: %d\n", d.Type)
	fmt.Println("start simulation")
	switch d.Type {
	case 0:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		ambient_sensor := NewAmbientSensorSimulator(client, d)
		ambient_sensor.ConnectAmbientSensor()
	case 1:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		air_conditioner := NewAirConditionerSimulator(client, d)
		air_conditioner.ConnectAirConditioner()
	case 2:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		air_conditioner := NewWashingMachineSimulator(client, d)
		air_conditioner.ConnectWashingMachine()
	case 3:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		lamp := NewLampSimulator(client, d)
		lamp.ConnectLamp()

	case 4:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		vehicleGate := NewVehicleGateSimulator(client, d)
		vehicleGate.ConnectVehicleGate()
	case 6:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		sp := NewSolarPanelSimulator(client, d)
		sp.ConnectSolarPanel()
	case 7:
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		sp := NewBatterySimulator(client, d)
		sp.ConnectBattery()
	default:
		fmt.Println("usao je u default case")
		fmt.Printf("Connecting device id=%d, Name=%s\n", d.ID, d.Name)
		lamp := NewLampSimulator(client, d)
		lamp.ConnectLamp()
	}
}

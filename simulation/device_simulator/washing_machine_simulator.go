package device_simulator

import (
	"encoding/json"
	"fmt"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicWMSwitch = "ac/switch/" // front salje sta se upalilo/ugasilo i ide do back-a
)

type WashingMachineSimulator struct {
	client mqtt.Client
	device models.WashingMachine
	off_on models.WMReceiveValue
}

func NewWashingMachineSimulator(client mqtt.Client, device models.Device) *WashingMachineSimulator {
	wm, err := config.GetWashingMachine(device.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	off_on := models.WMReceiveValue{}

	return &WashingMachineSimulator{
		client: client,
		device: wm,
		off_on: off_on,
	}
}

func (wm *WashingMachineSimulator) ConnectWashingMachine() {
	go SendHeartBeat(wm.client, wm.device.Device.Device.ID, wm.device.Device.Device.Name)
	config.SubscribeToTopic(wm.client, topicWMSwitch+strconv.Itoa(wm.device.Device.Device.ID), wm.HandleSwitchChange)
}

func (wm *WashingMachineSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	_, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	var washing_machine models.WMReceiveValue
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(msg.Payload()), &washing_machine)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// set values
	wm.off_on = washing_machine
	fmt.Println("PRIMLJENA PORUKA u ves masini")
	fmt.Println(wm.off_on)
}

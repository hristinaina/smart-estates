package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
	"strconv"
	"strings"
)

func (mc *MQTTClient) HandleVehicleApproached(_ mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	gate, err := mc.vehicleGateRepository.Get(id)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	licensePlate := string(msg.Payload())

	fmt.Printf("Car with license plate: %s, approached to gate: %s with id: %d\n", licensePlate,
		gate.ConsumptionDevice.Device.Name, gate.ConsumptionDevice.Device.Id)

	mc.CheckApproachedVehicle(gate, licensePlate)
}

func (mc *MQTTClient) CheckApproachedVehicle(gate models.VehicleGate, licensePlate string) {
	if gate.Mode == enumerations.Public {
		err := mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "open")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// TODO: publish that gate is open
		// TODO: add somewhere who is inside
	} else {

	}
}

package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
	"strconv"
	"strings"
	"time"
)

// TODO: how to save vehicles that have entered (to simulate exit)?
// send it to simulation and in simulation save and simulate exit
// TODO: if gate was already opened, no checking, just send who has entered (also implement manual opening gate on frontend)

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

	payload := string(msg.Payload())
	payloadTokens := strings.Split(payload, "+")
	licensePlate := payloadTokens[0]
	fmt.Printf("Payload: %s", payload)
	action := payloadTokens[1]

	fmt.Printf("Car with license plate: %s, approached to gate: %s with id: %d\n", licensePlate,
		gate.ConsumptionDevice.Device.Name, gate.ConsumptionDevice.Device.Id)

	if action == "enter" {
		mc.CheckApproachedVehicle(gate, licensePlate, "enter")
	} else {
		fmt.Printf("%s is leaving...", licensePlate)
		mc.CheckApproachedVehicle(gate, licensePlate, "exit")
	}
}

func (mc *MQTTClient) CheckApproachedVehicle(gate models.VehicleGate, licensePlate string, action string) {
	if !gate.IsOpen {
		if (gate.Mode == enumerations.Public) || (contains(gate.LicensePlates, licensePlate)) || (action == "exit") {
			_, err := mc.vehicleGateRepository.UpdateIsOpen(gate.ConsumptionDevice.Device.Id, true)
			if repositories.CheckIfError(err) {
				return
			}
			err = mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "open+"+licensePlate+"+"+action)
			if repositories.CheckIfError(err) {
				return
			}
			fmt.Printf("Published that someone has entered and that gate %s is open.\n",
				gate.ConsumptionDevice.Device.Name)
			select {
			case <-time.After(5 * time.Second):
				_, err = mc.vehicleGateRepository.UpdateIsOpen(gate.ConsumptionDevice.Device.Id, false)
				if repositories.CheckIfError(err) {
					return
				}
				err = mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "close+"+licensePlate)
				if repositories.CheckIfError(err) {
					return
				}
				fmt.Println("Car has entered. Gate is closing...")
			}
			// TODO: publish simulaciji da je vozilo uslo a u simulaciji sacekati rendom vrijeme i onda objaviti napustanje
		}
	} else {
		// TODO: publish that someone has entered and don't close gate after that
	}
}

func (mc *MQTTClient) Exit(gate models.VehicleGate, licensePlate string) {
	if !gate.IsOpen {

	}
}

func contains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

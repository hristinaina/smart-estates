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

type VehicleGate struct {
	VehicleGateId int    `json:"vehicle_gate_id"`
	LicensePlate  string `json:"license_plate"`
	Action        string `json:"action"`
	Success       bool   `json:"success"`
}

var (
	gateUpdateChan = make(chan VehicleGate)
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
		fmt.Printf("%s is leaving... %d\n", licensePlate, gate.ConsumptionDevice.Device.Id)
		mc.CheckApproachedVehicle(gate, licensePlate, "exit")
	}
}

// TODO: maybe create enum for action
func (mc *MQTTClient) CheckApproachedVehicle(gate models.VehicleGate, licensePlate string, action string) {
	if !gate.IsOpen {
		if (gate.Mode == enumerations.Public) || (contains(gate.LicensePlates, licensePlate)) || (action == "exit") {
			fmt.Println("usao")
			_, err := mc.vehicleGateRepository.UpdateIsOpen(gate.ConsumptionDevice.Device.Id, true)
			if repositories.CheckIfError(err) {
				return
			}
			err = mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "open+"+licensePlate+"+"+action)
			if repositories.CheckIfError(err) {
				return
			}
			go func() {
				select {
				case <-time.After(6 * time.Second):
					_, err = mc.vehicleGateRepository.UpdateIsOpen(gate.ConsumptionDevice.Device.Id, false)
					if repositories.CheckIfError(err) {
						return
					}
					err = mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "close+"+licensePlate)
					if repositories.CheckIfError(err) {
						return
					}
					mc.vehicleGateRepository.PostNewVehicleGateValue(gate, action, true, licensePlate)
					//setData(gate.ConsumptionDevice.Device.Id, licensePlate, action, true)
				}
			}()

		} else {
			mc.vehicleGateRepository.PostNewVehicleGateValue(gate, action, false, licensePlate)
			//setData(gate.ConsumptionDevice.Device.Id, licensePlate, action, false)

		}
	} else if (gate.Mode == enumerations.Public) || contains(gate.LicensePlates, licensePlate) || (action == "exit") {
		err := mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "open+"+licensePlate+"+"+action)
		if repositories.CheckIfError(err) {
			return
		}
		go func() {
			select {
			case <-time.After(6 * time.Second):
				err = mc.Publish(TopicVGOpenClose+strconv.Itoa(gate.ConsumptionDevice.Device.Id), "leave_open+"+licensePlate)
				if repositories.CheckIfError(err) {
					return
				}
				mc.vehicleGateRepository.PostNewVehicleGateValue(gate, action, true, licensePlate)
				//setData(gate.ConsumptionDevice.Device.Id, licensePlate, action, true)

			}
		}()
	}
	//mc.vehicleGateRepository.GetFromInfluxDb("-30m")
}

func contains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func setData(id int, licensePlate, action string, success bool) {
	gateUpdateChan <- VehicleGate{
		VehicleGateId: id,
		LicensePlate:  licensePlate,
		Action:        action,
		Success:       success,
	}
}

func GetNewGate() VehicleGate {
	return <-gateUpdateChan
}

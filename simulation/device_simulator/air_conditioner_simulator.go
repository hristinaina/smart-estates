package device_simulator

import (
	"encoding/json"
	"fmt"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicACSwitch = "ac/switch/"
	topicTemp     = "ac/temp"
)

type AirConditionerSimulator struct {
	client mqtt.Client
	device models.AirConditioner
	off_on models.ReceiveValue
}

func NewAirConditionerSimulator(client mqtt.Client, device models.Device) *AirConditionerSimulator {
	ac, err := config.GetAC(device.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	off_on := models.ReceiveValue{}

	return &AirConditionerSimulator{
		client: client,
		device: ac,
		off_on: off_on,
	}
}

func (ac *AirConditionerSimulator) ConnectAirConditioner() {
	go SendHeartBeat(ac.client, ac.device.Device.Device.ID, ac.device.Device.Device.Name)
	go ac.GenerateAirConditionerData()
	config.SubscribeToTopic(ac.client, topicACSwitch+strconv.Itoa(ac.device.Device.Device.ID), ac.HandleSwitchChange)
}

func (ac *AirConditionerSimulator) GenerateAirConditionerData() {
	temp := 20.0

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println(ac.off_on)
			if ac.off_on.Switch {
				fmt.Println("uslooo")
				switch ac.off_on.Mode {
				case "Heating":
					if temp < float64(ac.off_on.Temp) {
						if temp+0.5 > float64(ac.off_on.Temp) {
							temp = float64(ac.off_on.Temp)
						} else {
							temp += 0.5
						}
					} else {
						temp = float64(ac.off_on.Temp)
					}

				case "Cooling":
					if temp > float64(ac.off_on.Temp) {
						if temp-0.5 < float64(ac.off_on.Temp) {
							temp = float64(ac.off_on.Temp)
						} else {
							temp -= 0.5
						}
					} else {
						temp = float64(ac.off_on.Temp)
					}
				case "Automatic":
					if temp > float64(ac.off_on.Temp) {
						if temp-0.5 < float64(ac.off_on.Temp) {
							temp = float64(ac.off_on.Temp)
						} else {
							temp -= 0.5
						}
					} else if temp < float64(ac.off_on.Temp) {
						if temp+0.5 > float64(ac.off_on.Temp) {
							temp = float64(ac.off_on.Temp)
						} else {
							temp += 0.5
						}
					}
				case "Ventilation":
					// do not change temperature
				}
			} else {
				temp = ac.SendCurrentTemp()
			}
			// send on front
			data := map[string]interface{}{
				"id":   ac.device.Device.Device.ID,
				"temp": temp,
			}
			jsonString, err := json.Marshal(data)
			if err != nil {
				fmt.Println("greska")
			}
			config.PublishToTopic(ac.client, topicTemp, string(jsonString))

			fmt.Printf("Air Conditioner name=%s, id=%d, generated data: %f\n", ac.device.Device.Device.Name, ac.device.Device.Device.ID, temp)
		}
	}
}

func (ac *AirConditionerSimulator) SendCurrentTemp() float64 {
	openMeteoResponse, err := config.GetTemp()
	if err != nil {
		fmt.Printf("Error: %v \n", err.Error())
		return 20.0
	} else {
		temp := 0.8*openMeteoResponse.Current.Temperature2m + 15
		return temp
	}
}

func (ac *AirConditionerSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	_, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	var air_conditioner models.ReceiveValue
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(msg.Payload()), &air_conditioner)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// set values
	ac.off_on = air_conditioner

	// todo send values to back
	// ac.device.IsOn = sp.IsOn == true
	// fmt.Printf("Solar panel id=%d, switch status:\n", deviceId)
}

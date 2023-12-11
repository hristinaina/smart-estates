package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
	"strconv"
	"strings"
)

// TODO : create while loop that will do publishes

func (mc *MQTTClient) HandleValueChange(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("usaooooo")
	parts := strings.Split(msg.Topic(), "/")
	fmt.Println(msg.Topic())
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	device, err := mc.deviceRepository.Get(id)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// TODO: finish this
	fmt.Println("DEVICE LAMPPP: ", device.Id)

}

func (mc *MQTTClient) CheckValue(values []float32) {
	devices := mc.deviceRepository.GetAll()
	i := -1
	for _, device := range devices {
		i++
		if values[i] != device.LastValue {
			if device.Type == enumerations.Lamp {
				lamp, err := mc.lampRepository.Get(device.Id)
				if err != nil {
					fmt.Println("Error: ", err)
					// TODO: handle error
				} else {
					fmt.Println("Posting new lamp value to influxdb...")
					mc.PostNewLampValue(lamp)
				}
			}
			// TODO: handle other types
		}
	}
}

func (mc *MQTTClient) PostNewLampValue(lamp models.Lamp) {
	// TODO: implement this
}

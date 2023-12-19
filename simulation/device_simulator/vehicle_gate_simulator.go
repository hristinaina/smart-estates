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
	"net/http"
	"io/ioutil"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicApproached = "device/approached/"
	TopicVGOpenClose = "vg/open/"
)

type VehicleGateSimulator struct {
	client mqtt.Client
	device models.Device
	licensePlates []string
	consumptionDevice models.ConsumptionDevice
}

func NewVehicleGateSimulator(client mqtt.Client, device models.Device) *VehicleGateSimulator{
	parsedStrings := getAllLicensePlates()

	consumptionDevice, err := config.GetConsumptionDevice(device.ID)
	if err != nil {
		fmt.Println("Error while getting consumption device for lamp, id: " + strconv.Itoa(device.ID))
		return &VehicleGateSimulator {
			client: client,
			device: device,
			licensePlates: parsedStrings,
		}
	}

	return &VehicleGateSimulator {
		client: client,
		device: device,
		licensePlates: parsedStrings,
		consumptionDevice: consumptionDevice,
	}
}

func getAllLicensePlates() []string {
	url := "http://localhost:8081/api/vehicle-gate/license-plate"

	response, _ := http.Get(url)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	var parsedStrings []string

	err = json.Unmarshal([]byte(body), &parsedStrings)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return parsedStrings
}

func (sim *VehicleGateSimulator) ConnectVehicleGate() {
	go SendHeartBeat(sim.client, sim.device.ID, sim.device.Name)
	go sim.GenerateVehicleData()
	config.SubscribeToTopic(sim.client, TopicVGOpenClose+strconv.Itoa(sim.device.ID), sim.HandleLeaving)
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
	licensePlate := sim.getRandomLicensePlate()
	config.PublishToTopic(sim.client, config.TopicApproached+strconv.Itoa(sim.device.ID), licensePlate+"+enter")
	fmt.Println("Published to topic approached!")
}

func (sim *VehicleGateSimulator) HandleLeaving(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	payloadTokens := strings.Split(payload, "+")
	if (payloadTokens[0] == "open") {
		licensePlate := payloadTokens[1]
		action := payloadTokens[2]
		if action == "enter" {
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Float64()
			sec := int(randomNumber * 15)
			go func () {
				timerChan := time.After(time.Duration(sec) * time.Second)
				select {
				case <- timerChan:
					fmt.Printf("Left %s\n", licensePlate)
					config.PublishToTopic(sim.client, config.TopicApproached+strconv.Itoa(sim.device.ID), licensePlate+"+exit")
				}
			} ()
		}
	}
}

func (sim *VehicleGateSimulator) getRandomLicensePlate() string {
	x := len(sim.licensePlates)
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(x)
	return sim.licensePlates[index]
}
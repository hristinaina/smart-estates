package device_simulator

import (
	"fmt"
	// "math"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"simulation/config"
	"simulation/models"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicSwitch = "lamp/switch/"
)

type LampSimulator struct {
	switchOn          bool
	client            mqtt.Client
	device            models.Device
	consumptionDevice models.ConsumptionDevice
}

func NewLampSimulator(client mqtt.Client, device models.Device) *LampSimulator {
	//todo da se proslijedi samo deviceId (umjesto device) i posalje upit ka beku za dobavljane svih podataka za lampu
	// (jer device ima samo opste podatke)
	consumptionDevice, err := config.GetConsumptionDevice(device.ID)
	if err != nil {
		fmt.Println("Error while getting consumption device for lamp, id: " + strconv.Itoa(device.ID))
		return &LampSimulator{
			client:   client,
			device:   device,
			switchOn: false,
		}
	}

	return &LampSimulator{
		client:            client,
		device:            device,
		switchOn:          false,
		consumptionDevice: consumptionDevice,
	}
}

func (ls *LampSimulator) ConnectLamp() {
	go SendHeartBeat(ls.client, ls.device.ID, ls.device.Name)
	go ls.GenerateLampData()
	go ls.SendConsumption()
	config.SubscribeToTopic(ls.client, topicSwitch+strconv.Itoa(ls.device.ID), ls.HandleSwitchChange)
}

// todo get real device and it's consumption (only if its powering is netwok, if it is not then end the function)
// SendConsumprion Periodically send consumption
func (ls *LampSimulator) SendConsumption() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			scalingFactor := 1.0
			if ls.switchOn {
				scalingFactor = 0.8 + rand.Float64()*0.2 // get a number between 0.8 and 1.0
			} else {
				scalingFactor = 0.15 + rand.Float64()*0.2 // get a number between 0.15 and 0.35
			}
			consumed := ls.consumptionDevice.PowerConsumption * scalingFactor / 60 / 2 // divide by 60 and 2 to get consumption for previous 30s
			err := config.PublishToTopic(ls.client, config.TopicConsumption+strconv.Itoa(ls.device.ID), strconv.FormatFloat(consumed,
				'f', -1, 64))
			if err != nil {
				fmt.Printf("Error publishing message with the device: %s \n", ls.device.Name)
			} else {
				fmt.Printf("%s: Lamp with id=%d, Name=%s, consumed=%fkWh for previous 30s\n", time.Now().Format("15:04:05"),
					ls.device.ID, ls.device.Name, consumed)
			}
		}
	}
}

// GenerateLampData Simulate sending periodic Lamp data
func (ls *LampSimulator) GenerateLampData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if ls.switchOn {
				// Get the Unix timestamp from the current time
				// unixTimestamp := float64(time.Now().Unix())
				// sineValue := math.Sin(unixTimestamp)
				// percentage := math.Abs(math.Round(sineValue * 100))
				percentage := float64(ls.getOutsideBrightness())
				config.PublishToTopic(ls.client, config.TopicPayload+strconv.Itoa(ls.device.ID), strconv.FormatFloat(percentage,
					'f', -1, 64))
				fmt.Printf("Lamp name=%s, id=%d, generated data: %f\n", ls.device.Name, ls.device.ID, percentage)
			}
		}
	}
}

func (ls *LampSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	ls.switchOn = status == "true"
	fmt.Printf("Lamp id=%d, switch status: %s\n", deviceId, status)
}

func (ls *LampSimulator) getOutsideBrightness() int {
	hour := time.Now().Hour()

	if hour >= 6 && hour < 16 {
		return rand.Intn(30) + 70
	} else {
		return rand.Intn(30)
	}
}

package device_simulator

import (
	"fmt"
	"math"
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
	switchOn bool
	client   mqtt.Client
	device   models.Device
}

func NewLampSimulator(client mqtt.Client, device models.Device) *LampSimulator {
	//todo da se proslijedi samo deviceId (umjesto device) i posalje upit ka beku za dobavljane svih podataka za lampu
	// (jer device ima samo opste podatke)
	return &LampSimulator{
		client:   client,
		device:   device,
		switchOn: false,
	}
}

func (ls *LampSimulator) ConnectLamp() {
	go SendHeartBeat(ls.client, ls.device)
	go ls.GenerateLampData()
	config.SubscribeToTopic(ls.client, topicSwitch+strconv.Itoa(ls.device.ID), ls.HandleSwitchChange)
}

// GenerateLampData Simulate sending periodic Lamp data
func (ls LampSimulator) GenerateLampData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ls.switchOn {
				// Get the Unix timestamp from the current time
				unixTimestamp := float64(time.Now().Unix())
				sineValue := math.Sin(unixTimestamp)
				percentage := math.Abs(math.Round(sineValue * 100))
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

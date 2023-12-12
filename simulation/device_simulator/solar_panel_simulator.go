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
	topicSPSwitch = "sp/switch/"
)

type SolarPanelSimulator struct {
	client mqtt.Client
	device models.SolarPanel
}

func NewSolarPanelSimulator(client mqtt.Client, device models.Device) *SolarPanelSimulator {
	sp, err := config.GetSP(device.ID)
	if err != nil {
		return nil
	}
	return &SolarPanelSimulator{
		client: client,
		device: sp,
	}
}

func (ls *SolarPanelSimulator) ConnectSolarPanel() {
	go SendHeartBeat(ls.client, ls.device.Device.ID, ls.device.Device.Name)
	go ls.GenerateSolarPanelData()
	config.SubscribeToTopic(ls.client, topicSwitch+strconv.Itoa(ls.device.Device.ID), ls.HandleSwitchChange)
}

func (ls *SolarPanelSimulator) GenerateSolarPanelData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ls.device.IsOn {
				// Get the Unix timestamp from the current time
				unixTimestamp := float64(time.Now().Unix())
				sineValue := math.Sin(unixTimestamp)
				percentage := math.Abs(math.Round(sineValue * 100))
				config.PublishToTopic(ls.client, config.TopicPayload+strconv.Itoa(ls.device.Device.ID), strconv.FormatFloat(percentage,
					'f', -1, 64))
				fmt.Printf("Solar panel name=%s, id=%d, generated data: %f\n", ls.device.Device.Name, ls.device.Device.ID, percentage)
			}
		}
	}
}

// todo da ide info sa beka?
func (ls *SolarPanelSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	ls.device.IsOn = status == "true"
	fmt.Printf("Solar panel id=%d, switch status: %s\n", deviceId, status)
}

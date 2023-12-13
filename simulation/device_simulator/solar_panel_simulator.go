package device_simulator

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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
	config.SubscribeToTopic(ls.client, topicSPSwitch+strconv.Itoa(ls.device.Device.ID), ls.HandleSwitchChange)
}

func (ls *SolarPanelSimulator) GenerateSolarPanelData() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ls.device.IsOn {
				// SolarRadiation is in W/m^2
				openMeteoResponse, err := config.GetSolarRadiation(45.45, 19) //todo get real lat and long
				if err != nil {
					fmt.Printf("Error: %v \n", err.Error())
				} else {
					solarRadiation := openMeteoResponse.Hourly.DirectNormalIrradiance[time.Now().Hour()]
					fmt.Println(solarRadiation)
					electricityProduction := calculateEletricityProduction(ls.device, solarRadiation)
					config.PublishToTopic(ls.client, config.TopicPayload+strconv.Itoa(ls.device.Device.ID), strconv.FormatFloat(electricityProduction,
						'f', -1, 64))
					fmt.Printf("Solar panel name=%s, id=%d, generated data: %f\n", ls.device.Device.Name, ls.device.Device.ID, electricityProduction)
				}

			}
		}
	}
}

// solar radiation is used to scale electricity depending on sun intensity
func calculateEletricityProduction(sp models.SolarPanel, radiation float64) float64 {
	if radiation == 0 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	hourProduction := sp.SurfaceArea * sp.Efficiency * 10 // multiply by 1000 and divide by 100 (percentage)
	minuteProduction := hourProduction / 60
	scalingFactor := radiation/(radiation+2) + rand.Float64()*(0.25-0.2) + 0.2 // create scaling factor depending on sun radiation and add random factor to show changes
	electricity := math.Round(minuteProduction*scalingFactor*1e2) / 1e2        // scale eletricity with scaling factor (shown in Wh) and rounds to 2 decimal places
	return electricity
}

func (ls *SolarPanelSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	var sp models.SolarPanel
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(msg.Payload()), &sp)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	fmt.Println(sp)
	ls.device.IsOn = sp.IsOn == true
	fmt.Printf("Solar panel id=%d, switch status: %t\n", deviceId, sp.IsOn)
}

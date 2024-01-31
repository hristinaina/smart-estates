package device_simulator

import (
	"simulation/config"
	models "simulation/models"
	"strconv"
	"strings"
	"time"
	"fmt"
	"net/http"
	// "io/ioutil"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	// TODO: add topics here
	// topicApproached = "device/approached/"
	// TopicVGOpenClose = "vg/open/"
	TurnSprinklerON = "sprinkler/on"
)

type SprinklerSimulator struct {
	client mqtt.Client
	device models.Device
	consumptionDevice models.ConsumptionDevice
}

func NewSprinklerSimulator(client mqtt.Client, device models.Device) *SprinklerSimulator {
	consumptionDevice, err := config.GetConsumptionDevice(device.ID)
	if err != nil {
		fmt.Println("Error while getting consumption device for lamp, id: " + strconv.Itoa(device.ID))
		return &SprinklerSimulator {
			client: client,
			device: device,
		}
	}
	return &SprinklerSimulator {
		client: client,
		device: device,
		consumptionDevice: consumptionDevice,
	}
}

func (sim *SprinklerSimulator) ConnectSprinkler() {
	go SendHeartBeat(sim.client, sim.device.ID, sim.device.Name)
	go sim.CheckScheduledModes()
	// config.SubscribeToTopic(sim.client, TopicVGOpenClose+strconv.Itoa(sim.device.ID), sim.HandleLeaving)
}

func (sim *SprinklerSimulator) CheckScheduledModes() {
	// TODO: ovdje while petlja koja svakih 60 sekundi dobavlja zakazane termine i provjerava da li treba upaliti prskalicu
	// ukoliko je treba upaliti, objavi na topic koji ce hvatati front i back
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scheduledModes := sim.getSheduledModes();
			mode := sim.checkScheduledMode(scheduledModes);
			fmt.Println("MODEEEEEE")
			fmt.Println(mode)
			if mode.Id != 0 {
				fmt.Println("PUBLISHINGGGG")
				err := config.PublishToTopic(sim.client, config.TurnSprinklerON+strconv.Itoa(sim.device.ID), strconv.Itoa(mode.Id))
				if err != nil {
					fmt.Printf("Error publishing message with the device: %s \n", sim.device.Name)
				} else {
					fmt.Println("Successfully publishedddd")
				}
			}
		}
	}
}

func (sim *SprinklerSimulator) getSheduledModes() []models.SprinklerSpecialMode {
	url := "http://localhost:8081/api/sprinkler/mode/" + strconv.Itoa(sim.device.ID);

		// Make the HTTP request
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return nil
		}
		defer response.Body.Close()
	
		// Check if the response status code is OK
		if response.StatusCode != http.StatusOK {
			fmt.Println("HTTP request failed with status code:", response.StatusCode)
			return nil
		}
	
		// Decode the JSON response into a slice of SprinklerSpecialMode objects
		var modes []models.SprinklerSpecialMode
		if err := json.NewDecoder(response.Body).Decode(&modes); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil
		}
	
		// Print the parsed data
		fmt.Printf("Parsed SprinklerSpecialModes: %+v\n", modes)
		return modes
}

func (sim *SprinklerSimulator) checkScheduledMode(modes []models.SprinklerSpecialMode) models.SprinklerSpecialMode {
	currentTime := time.Now()
	currentDay := currentTime.Weekday()

	for _, s := range modes {
		one := sim.isTimeInRange(currentTime, sim.parseTime(s.StartTime), sim.parseTime(s.EndTime))
		if one && sim.isDayInSchedule(currentDay, s.SelectedDays) {
			return s
		}
	}
	return models.SprinklerSpecialMode{}
}

func (sim *SprinklerSimulator) isTimeInRange(current, start, end time.Time) bool {
	currentTime := current.Hour()*60 + current.Minute()
	startTime := start.Hour()*60 + start.Minute()
	endTime := end.Hour()*60 + end.Minute()

	return (currentTime >= startTime && currentTime <= endTime) || (currentTime <= startTime && currentTime <= endTime)
}

func (sim *SprinklerSimulator) isDayInSchedule(currentDay time.Weekday, days string) bool {
	if days == "" {
		return true
	}
	if strings.Contains(days, currentDay.String()) {
		return true
	} else {
		return false
	}
}

func (sim *SprinklerSimulator) parseTime(timeString string) time.Time {
	layout := "15:04:05"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		panic(err)
	}
	return t
}

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
	TurnSprinklerON = "sprinkler/on/"
	TurnSprinklerOFF = "sprinkler/off/"
)

type SprinklerSimulator struct {
	isOn bool
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
			isOn: false,
		}
	}
	return &SprinklerSimulator {
		client: client,
		device: device,
		consumptionDevice: consumptionDevice,
		isOn: false,
	}
}

func (sim *SprinklerSimulator) ConnectSprinkler() {
	go SendHeartBeat(sim.client, sim.device.ID, sim.device.Name)
	go sim.CheckScheduledModes()
}

func (sim *SprinklerSimulator) CheckScheduledModes() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scheduledModes := sim.getSheduledModes();
			mode := sim.checkScheduledMode(scheduledModes);
			if (mode.Id != 0) && (sim.isOn == false) {
				err := config.PublishToTopic(sim.client, config.TurnSprinklerON+strconv.Itoa(sim.device.ID), strconv.Itoa(mode.Id))
				sim.isOn = true
				if err != nil {
					fmt.Printf("Error publishing message with the device: %s \n", sim.device.Name)
				} else {
					fmt.Println(config.TurnSprinklerON+strconv.Itoa(sim.device.ID))
				}
			} else {
				if sim.isCurrentTimeWithin70SecondsAfterEndTime(scheduledModes) {
					if sim.isOn == true {
						err := config.PublishToTopic(sim.client, config.TurnSprinklerOFF+strconv.Itoa(sim.device.ID), strconv.Itoa(mode.Id))
						sim.isOn = false
						if err != nil {
							fmt.Printf("Error publishing message with the device: %s \n", sim.device.Name)
						} else {
							fmt.Println(config.TurnSprinklerOFF+strconv.Itoa(sim.device.ID))
						}
					}
					
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

	// Check if the range spans across two different days
	if startTime > endTime {
		return currentTime >= startTime || currentTime <= endTime
	}

	// Check if the current time is within the range
	return currentTime >= startTime && currentTime <= endTime
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

func (sim *SprinklerSimulator) isCurrentTimeWithin70SecondsAfterEndTime(modes []models.SprinklerSpecialMode) bool {
    current := time.Now()
	currentDay := current.Weekday()

	for _, s := range modes {
		if sim.isDayInSchedule(currentDay, s.SelectedDays) {
			end := sim.parseTime(s.EndTime)
			end = time.Date(current.Year(), current.Month(), current.Day(), end.Hour(), end.Minute(), end.Second(), current.Nanosecond(), current.Location())
			diff := current.Sub(end).Seconds()
			// TODO: change this to 100 + smth later (between 60 and 120)
			if (diff <= 119) && (diff >= 0) {
				return true
			}
		}
	}
	return false
}

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
	topicACSwitch = "ac/switch/" // front salje sta se upalilo/ugasilo i ide do back-a
	topicTemp     = "ac/temp"    // salje temp na front
	topicAction   = "ac/action"  // slanje na front i back zakazan termin
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
	go ac.SendScheduledData()
	go ac.SendConsumption()
	config.SubscribeToTopic(ac.client, topicACSwitch+strconv.Itoa(ac.device.Device.Device.ID), ac.HandleSwitchChange)
}

func (ls *AirConditionerSimulator) SendConsumption() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			scalingFactor := 0.8 + rand.Float64()*0.2                              // get a number between 0.8 and 1.0
			consumed := ls.device.Device.PowerConsumption * scalingFactor / 60 / 2 // divide by 60 and 2 to get consumption for previous 30s
			err := config.PublishToTopic(ls.client, config.TopicConsumption+strconv.Itoa(ls.device.Device.Device.ID), strconv.FormatFloat(consumed,
				'f', -1, 64))
			if err != nil {
				fmt.Printf("Error publishing message with the device: %s \n", ls.device.Device.Device.Name)
			} else {
				fmt.Printf("%s: Lamp with id=%d, Name=%s, consumed=%fkWh for previous 30s\n", time.Now().Format("15:04:05"),
					ls.device.Device.Device.ID, ls.device.Device.Device.Name, consumed)
			}
		}
	}
}

func (ac *AirConditionerSimulator) GenerateAirConditionerData() {
	currentTemp := ac.GetCurrentTemp()
	temp := currentTemp

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ac.off_on.Switch {
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
				if temp < currentTemp {
					if temp+0.5 > currentTemp {
						temp = currentTemp
					} else {
						temp += 0.5
					}
				} else if temp > currentTemp {
					if temp-0.5 < currentTemp {
						temp = currentTemp
					} else {
						temp -= 0.5
					}
				}
			}
			// send on front
			data := map[string]interface{}{
				"id":   ac.device.Device.Device.ID,
				"temp": math.Round(temp*100) / 100,
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

func (ac *AirConditionerSimulator) GetCurrentTemp() float64 {
	openMeteoResponse, err := config.GetTemp()
	if err != nil {
		fmt.Printf("Error: %v \n", err.Error())
		return 20.0
	} else {
		temp := 0.5*openMeteoResponse.Current.Temperature2m + 15
		return math.Round(temp*100) / 100
	}
}

func (ac *AirConditionerSimulator) SendScheduledData() {
	// temp := ac.GetCurrentTemp()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if ac.device.SpecialMode[0].StartTime != "" {
				s := ac.checkScheduledMode()

				if s.StartTime != "" {
					// send on front
					data := map[string]interface{}{
						"id":     ac.device.Device.Device.ID,
						"mode":   s.Mode,
						"switch": true,
						"temp":   s.Temperature,
					}
					jsonString, err := json.Marshal(data)
					if err != nil {
						fmt.Println("greska")
					}
					config.PublishToTopic(ac.client, topicAction, string(jsonString))

					fmt.Printf("Turn on: %s \n", s.Mode)
				}
			}
		}
	}
}

func (ac *AirConditionerSimulator) checkScheduledMode() models.SpecialMode {
	currentTime := time.Now()
	currentDay := currentTime.Weekday()

	for _, s := range ac.device.SpecialMode {
		one := isTimeInRange(currentTime, parseTime(s.StartTime), parseTime(s.EndTime))
		if one && isDayInSchedule(currentDay, s.SelectedDays) {
			fmt.Printf("Schould turn on: %s, temp: %f\n", s.Mode, s.Temperature)
			return s
		}
	}
	return models.SpecialMode{}
}

func isTimeInRange(current, start, end time.Time) bool {
	currentTime := current.Hour()*60 + current.Minute()
	startTime := start.Hour()*60 + start.Minute()
	endTime := end.Hour()*60 + end.Minute()

	return (currentTime >= startTime && currentTime <= endTime) || (currentTime <= startTime && currentTime <= endTime)
}

func isDayInSchedule(currentDay time.Weekday, days string) bool {
	if strings.Contains(days, currentDay.String()) {
		return true
	} else {
		return false
	}
}

func parseTime(timeString string) time.Time {
	layout := "15:04:05"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		panic(err)
	}
	return t
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
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(ac.off_on)
}

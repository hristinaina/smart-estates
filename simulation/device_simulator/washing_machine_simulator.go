package device_simulator

import (
	"encoding/json"
	"fmt"
	"simulation/config"
	"simulation/models"
	"sort"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topicWMSwitch    = "wm/switch/"  // front salje sta se upalilo/ugasilo i ide do back-a uopste ne ide do simulacije
	topicScheduled   = "wm/schedule" // slanje na front da se upali zakazan rezim
	topicGetSchedule = "wm/get/"     // prima sa fronta da li je nesto novo zakazano
)

type WashingMachineSimulator struct {
	client       mqtt.Client
	device       models.WashingMachine
	off_on       models.WMReceiveValue
	stopSchedule chan struct{}
}

func NewWashingMachineSimulator(client mqtt.Client, device models.Device) *WashingMachineSimulator {
	wm, err := config.GetWashingMachine(device.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	off_on := models.WMReceiveValue{}

	return &WashingMachineSimulator{
		client: client,
		device: wm,
		off_on: off_on,
	}
}

func (wm *WashingMachineSimulator) ConnectWashingMachine() {
	go SendHeartBeat(wm.client, wm.device.Device.Device.ID, wm.device.Device.Device.Name)
	go wm.ScheduleMode()
	config.SubscribeToTopic(wm.client, topicGetSchedule+strconv.Itoa(wm.device.Device.Device.ID), wm.HandleSwitchChange)
}

func (wm *WashingMachineSimulator) ScheduleMode() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	first := wm.getFirstScheduledToday()
	fmt.Println(first)

	loc, _ := time.LoadLocation("Europe/Belgrade")

	if first != nil {
		for {
			select {
			case <-ticker.C:
				now := time.Now().In(loc)

				startTime, err := time.ParseInLocation("2006-01-02 15:04:05", first.StartTime, loc)
				if err != nil {
					fmt.Printf("Error parsing start time: %v\n", err)
					continue
				}

				if now.After(startTime) || now.Equal(startTime) {
					fmt.Println("Time to execute scheduled mode!")
					data := map[string]interface{}{
						"id":       wm.device.Device.Device.ID,
						"mode":     first.ModeId,
						"switchOn": true,
					}
					jsonString, err := json.Marshal(data)
					if err != nil {
						fmt.Println(err)
					}
					config.PublishToTopic(wm.client, topicScheduled, string(jsonString))
					return
				}

			case <-wm.stopSchedule:
				fmt.Println("ScheduleMode interrupted")
				return
			}
		}
	}
}

func (wm *WashingMachineSimulator) getFirstScheduledToday() *models.ScheduledMode {
	scheduledMode, err := config.GetWashingMachineScheduledMode(wm.device.Device.Device.ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	wm.device.ScheduledMode = scheduledMode

	today := time.Now().Format("2006-01-02")

	// Sortirajte listu zakazanih termina po vremenu početka
	sort.Slice(wm.device.ScheduledMode, func(i, j int) bool {
		startTimeI, _ := time.Parse("2006-01-02 15:04:05", wm.device.ScheduledMode[i].StartTime)
		startTimeJ, _ := time.Parse("2006-01-02 15:04:05", wm.device.ScheduledMode[j].StartTime)
		return startTimeI.Before(startTimeJ)
	})

	// Prolazite kroz sortiranu listu i pronađite termin koji treba da se izvrši danas
	for _, term := range wm.device.ScheduledMode {
		startTime, err := time.Parse("2006-01-02 15:04:05", term.StartTime)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			continue
		}

		if startTime.Format("2006-01-02") == today {
			return &term
		}
	}

	return nil
}

func (wm *WashingMachineSimulator) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	_, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	var washing_machine models.WMReceiveValue
	err = json.Unmarshal([]byte(msg.Payload()), &washing_machine)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	if washing_machine.Get {
		fmt.Println("PRIMLJENA PORUKA u ves masini")

		if wm.stopSchedule != nil {
			close(wm.stopSchedule)
		}
		wm.stopSchedule = make(chan struct{})
		go wm.ScheduleMode()
	}

}

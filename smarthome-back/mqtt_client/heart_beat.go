package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func (mc *MQTTClient) HandleHeartBeat(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	device, err := mc.deviceRepository.Get(deviceId)
	if !device.IsOnline {
		err := mc.Publish(TopicStatusChanged+strconv.Itoa(deviceId), "online")
		if err != nil {
			fmt.Println(err)
		}
		//todo save to influxdb
	}
	device.IsOnline = true
	device.StatusTimeStamp = mysql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	mc.deviceRepository.Update(device)
	fmt.Printf("Device is online, id=%d\n", deviceId)
}

// CheckDeviceStatus function that checks if there is a device that has disconnected
func (mc *MQTTClient) CheckDeviceStatus() {
	offlineTimeout := 30 * time.Second
	devices := mc.deviceRepository.GetAll()
	for _, device := range devices {
		if device.IsOnline && time.Since(device.StatusTimeStamp.Time) > offlineTimeout {
			fmt.Printf("Device with id=%d is offline.\n", device.Id)
			device.IsOnline = false
			device.StatusTimeStamp = mysql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
			mc.deviceRepository.Update(device)
			err := mc.Publish(TopicStatusChanged+strconv.Itoa(device.Id), "offline")
			//todo save to influxdb
			if err != nil {
				return
			}
		}
	}
}

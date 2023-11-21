package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
	"time"
)

var heartBeats = make(map[int]time.Time)

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func HandleHeartBeat(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	if status == "offline" {
		return
	}
	heartBeats[deviceId] = time.Now()
	fmt.Printf("Device is online, id=%d\n", deviceId)
}

// CheckDeviceStatus function that checks if there is a device that has disconnected
func CheckDeviceStatus(client mqtt.Client) {
	offlineTimeout := 30 * time.Second

	for deviceID, timestamp := range heartBeats {
		if time.Since(timestamp) > offlineTimeout {
			fmt.Printf("Device with id=%d is offline.\n", deviceID)
			delete(heartBeats, deviceID)
			err := PublishToTopic(client, TopicOnline+strconv.Itoa(deviceID), "offline")
			if err != nil {
				return
			}
		}
	}
}

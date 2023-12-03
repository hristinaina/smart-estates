package mqtt_client

import (
	"database/sql"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
	"time"
)

var heartBeats = make(map[int]time.Time)

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func (mc *MQTTClient) HandleHeartBeat(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	_, ok := heartBeats[deviceId]
	if !ok {
		saveToDb(mc.db, deviceId, true)
		err := mc.Publish(TopicStatusChanged+strconv.Itoa(deviceId), "online")
		if err != nil {
			fmt.Println(err)
		}
	}
	heartBeats[deviceId] = time.Now()
	fmt.Printf("Device is online, id=%d\n", deviceId)
}

// CheckDeviceStatus function that checks if there is a device that has disconnected
func (mc *MQTTClient) CheckDeviceStatus() {
	offlineTimeout := 30 * time.Second

	for deviceID, timestamp := range heartBeats {
		if time.Since(timestamp) > offlineTimeout {
			fmt.Printf("Device with id=%d is offline.\n", deviceID)
			delete(heartBeats, deviceID)
			saveToDb(mc.db, deviceID, false)
			err := mc.Publish(TopicStatusChanged+strconv.Itoa(deviceID), "offline")
			if err != nil {
				return
			}
		}
	}
}

func saveToDb(db *sql.DB, deviceID int, flag bool) {
	query := "UPDATE device SET IsOnline = ? WHERE ID = ?"
	_, err := db.Exec(query, flag, deviceID)
	if err != nil {
		fmt.Println("Failed to update device status")
	}
}

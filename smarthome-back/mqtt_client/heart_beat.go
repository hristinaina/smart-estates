package mqtt_client

import (
	"fmt"
	models "smarthome-back/models/devices"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
		device.IsOnline = true
		device.StatusTimeStamp = mysql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		saveToInfluxDb(mc.influxDb, device)
	} else {
		device.IsOnline = true
		device.StatusTimeStamp = mysql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	mc.deviceRepository.Update(device)
	fmt.Printf("Device is online, id=%d\n", deviceId)
}

func (mc *MQTTClient) StartDeviceStatusThread() {
	// Periodically check if the device is still online
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// This block will be executed every time the ticker ticks
				fmt.Println("checking device status...")
				mc.checkDeviceStatus()
			}
		}
	}()
}

// CheckDeviceStatus function that checks if there is a device that has disconnected
func (mc *MQTTClient) checkDeviceStatus() {
	offlineTimeout := 60 * time.Second
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
			saveToInfluxDb(mc.influxDb, device)
			if err != nil {
				return
			}
		}
	}
}

func saveToInfluxDb(client influxdb2.Client, device models.Device) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	p := influxdb2.NewPoint("device_status", //table
		map[string]string{"device_id": strconv.Itoa(device.Id)}, //tag
		map[string]interface{}{"status": func() int {
			if device.IsOnline {
				return 1
			} else {
				return 0
			}
		}()}, //field
		device.StatusTimeStamp.Time)

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved status change to influxdb")
}

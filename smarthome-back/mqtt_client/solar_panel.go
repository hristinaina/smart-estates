package mqtt_client

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dto"
	models "smarthome-back/models/devices"
	"strconv"
	"strings"
	"time"
)

func (mc *MQTTClient) HandleSPSwitch(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	device := mc.solarPanelRepository.Get(deviceId)
	// Unmarshal the JSON string into the struct
	var data dto.DeviceDTO
	err = json.Unmarshal([]byte(msg.Payload()), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	device.IsOn = data.IsOn == true
	mc.solarPanelRepository.UpdateSP(device)
	//mc.deviceRepository.Update(device.Device)
	saveSwitchChangeToInfluxDb(mc.influxDb, device, data.UserId)
	fmt.Printf("Solar panel: name=%s, id=%d, changed switch to %t \n", device.Device.Name, device.Device.Id, device.IsOn)
}

func saveSwitchChangeToInfluxDb(client influxdb2.Client, device models.SolarPanel, id int) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("solar_panel", //table
		map[string]string{"device_id": strconv.Itoa(device.Device.Id), "user_id": strconv.Itoa(id)}, //tag
		map[string]interface{}{"isOn": func() int {
			if device.IsOn {
				return 1
			} else {
				return 0
			}
		}()}, //field
		time.Now())

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved sp switch change to influxdb")
}

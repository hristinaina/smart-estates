package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	status := string(msg.Payload())
	device.IsOn = status == "true"
	mc.solarPanelRepository.UpdateSP(device)
	//mc.deviceRepository.Update(device.Device)
	saveSwitchChangeToInfluxDb(mc.influxDb, device)
	fmt.Printf("Solar panel: name=%s, id=%d, changed switch to %t \n", device.Device.Name, device.Device.Id, device.IsOn)
}

func saveSwitchChangeToInfluxDb(client influxdb2.Client, device models.SolarPanel) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	//todo dodati korisnika koji je izvrsio akciju? Sa fronta dobavim userId i mogu slati kao json string
	p := influxdb2.NewPoint("solar_panel", //table
		map[string]string{"device_id": strconv.Itoa(device.Device.Id)}, //tag
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

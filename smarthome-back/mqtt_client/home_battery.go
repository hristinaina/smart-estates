package mqtt_client

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	models "smarthome-back/models/devices"
	"strconv"
	"time"
)

func (mc *MQTTClient) StartConsumptionThread() {
	// Periodically check devices consumption
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// This block will be executed every time the ticker ticks
				fmt.Println("checking device consumption...")
				mc.handleConsumption()
			}
		}
	}()
}

func (mc *MQTTClient) handleConsumption() {
	//todo dobavi rs
	realEstates, err := mc.realEstateRepository.GetAll()
	if err != nil {
		return
	}
	for _, value := range realEstates {
		fmt.Println("rs id:", value.Id)
		devices, err := mc.deviceRepository.GetConsumptionDevicesByEstateId(value.Id)
		if err != nil {
			return
		}
		for _, device := range devices {
			fmt.Println("\tdevice id", device.Device.Id)
		}
	}
	//todo dobavi consumption devices za taj rs
	//todo dobavi baterije za taj rs

}

func saveConsumptionToInfluxDb(client influxdb2.Client, device models.SolarPanel, email string) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("solar_panel", //table
		map[string]string{"device_id": strconv.Itoa(device.Device.Id), "user_id": email}, //tag
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

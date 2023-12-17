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

func (mc *MQTTClient) HandleConsumption(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the JSON string into the struct
	valueStr := string(msg.Payload())
	consumptionValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return
	}

	realEstate, err := mc.realEstateRepository.Get(deviceId)
	if err != nil {
		return
	}

	batteries, err := mc.homeBatteryRepository.GetAllByEstateId(realEstate.Id)
	if err != nil {
		return
	}
	fmt.Println("deviceID: ", deviceId, " realEstateID: ", realEstate.Id)
	for _, hb := range batteries {
		fmt.Println("\tbattery id", hb.Device.Id)
		if hb.CurrentValue-consumptionValue >= 0 { //end
			hb.CurrentValue = hb.CurrentValue - consumptionValue
			consumptionValue = 0
			//todo save battery to db and info to influx
		} else { //continue
			consumed := hb.CurrentValue
			consumptionValue = consumptionValue - hb.CurrentValue
			hb.CurrentValue = 0
			fmt.Println(consumed)
			//todo save battery to db and info(consumed) to influx
		}
	}
	if consumptionValue != 0 {
		//todo electrodistribution (influxdb)
	}

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

package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

	device, err := mc.deviceRepository.Get(deviceId)
	if err != nil {
		return
	}
	batteries, err := mc.homeBatteryRepository.GetAllByEstateId(device.RealEstate)
	if err != nil {
		return
	}
	fmt.Println("deviceID: ", deviceId, " realEstateID: ", device.RealEstate)
	for _, hb := range batteries {
		fmt.Println("\tbattery id", hb.Device.Id)
		if hb.CurrentValue-consumptionValue >= 0 { //end
			hb.CurrentValue = hb.CurrentValue - consumptionValue
			saveConsumptionToInfluxDb(mc.influxDb, device.RealEstate, device.Id, strconv.Itoa(hb.Device.Id), consumptionValue)
			mc.homeBatteryRepository.Update(hb)
			SaveHBDataToInfluxDb(mc.influxDb, hb.Device.Id, hb.CurrentValue)
			consumptionValue = 0
			break
		} else { //continue
			consumed := hb.CurrentValue
			consumptionValue = consumptionValue - hb.CurrentValue
			hb.CurrentValue = 0
			saveConsumptionToInfluxDb(mc.influxDb, device.RealEstate, device.Id, strconv.Itoa(hb.Device.Id), consumed)
			mc.homeBatteryRepository.Update(hb)
			SaveHBDataToInfluxDb(mc.influxDb, hb.Device.Id, hb.CurrentValue)
		}
	}
	if consumptionValue != 0 {
		saveConsumptionToInfluxDb(mc.influxDb, device.RealEstate, device.Id, "electrical_distribution", consumptionValue)
	}

}

func saveConsumptionToInfluxDb(client influxdb2.Client, estateId, deviceId int, batteryId string, electricity float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("consumption", //table
		map[string]string{"device_id": strconv.Itoa(deviceId), "estate_id": strconv.Itoa(estateId), "battery_id": batteryId},
		map[string]interface{}{"electricity": electricity},
		time.Now())

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved consumption to influxdb")
}

func SaveHBDataToInfluxDb(client influxdb2.Client, batteryId int, currentValue float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("home_battery", //table
		map[string]string{"device_id": strconv.Itoa(batteryId)},
		map[string]interface{}{"value": strconv.FormatFloat(currentValue, 'f', -1, 64)},
		time.Now())

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved consumption to influxdb")
}

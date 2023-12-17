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

func (mc *MQTTClient) HandleHBData(client mqtt.Client, msg mqtt.Message) {
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
	mc.handleConsumption(device, consumptionValue)
	fmt.Printf("Device: name=%s, id=%d, consumption value %f \n", device.Name, device.Id, consumptionValue)
}

func (mc *MQTTClient) handleConsumption(device models.Device, consumptionValue float64) {
	batteries, err := mc.homeBatteryRepository.GetAllByEstateId(device.RealEstate)
	if err != nil {
		return
	}
	surplus := 0.0
	if len(batteries) != 0 {
		valuePerBattery := consumptionValue / float64(len(batteries))
		// value is divided between batteries and each battery takes the same value
		surplus = mc.calculateConsumptionForBatteries(batteries, device.RealEstate, device.Id, valuePerBattery, false)
		if surplus != 0.0 {
			// surplus=what was left (if one of the batteries was full) is sent to batteries again (not divided)
			surplus = mc.calculateConsumptionForBatteries(batteries, device.RealEstate, device.Id, surplus, true)
		}
	} else if len(batteries) == 0 {
		saveConsumptionToInfluxDb(mc.influxDb, device.RealEstate, device.Id, "electrical_distribution", consumptionValue)
	} else if surplus != 0.0 {
		saveConsumptionToInfluxDb(mc.influxDb, device.RealEstate, device.Id, "electrical_distribution", surplus)
	}
}

func (mc *MQTTClient) calculateConsumptionForBatteries(batteries []models.HomeBattery, realEstateId int, deviceId int, consumptionValue float64, isSurplus bool) float64 {
	surplus := 0.0
	for _, hb := range batteries {
		if !hb.Device.IsOnline {
			if !isSurplus {
				surplus = surplus + consumptionValue
			}
			continue
		}
		if hb.CurrentValue-consumptionValue >= 0 { //end
			//everything that was supposed to be taken from battery was successfully taken
			hb.CurrentValue = hb.CurrentValue - consumptionValue
			saveConsumptionToInfluxDb(mc.influxDb, realEstateId, deviceId, strconv.Itoa(hb.Device.Id), consumptionValue)
			mc.homeBatteryRepository.Update(hb)
			SaveHBDataToInfluxDb(mc.influxDb, hb.Device.Id, hb.CurrentValue)
			if isSurplus {
				consumptionValue = 0
				break
			}
		} else { //continue
			//not everything that was supposed to be taken from the battery was taken
			consumed := hb.CurrentValue
			hb.CurrentValue = 0
			if isSurplus {
				consumptionValue = consumptionValue - hb.CurrentValue
			} else {
				surplus = surplus + consumptionValue - consumed
			}
			saveConsumptionToInfluxDb(mc.influxDb, realEstateId, deviceId, strconv.Itoa(hb.Device.Id), consumed)
			mc.homeBatteryRepository.Update(hb)
			SaveHBDataToInfluxDb(mc.influxDb, hb.Device.Id, hb.CurrentValue)
		}
	}
	if isSurplus {
		return consumptionValue
	} else {
		return surplus
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

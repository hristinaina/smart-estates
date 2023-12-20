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
	var data dto.SolarPanelDTO
	err = json.Unmarshal([]byte(msg.Payload()), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	device.IsOn = data.IsOn == true
	mc.solarPanelRepository.UpdateSP(device)
	//mc.deviceRepository.Update(device.Device)
	saveSPSwitchChangeToInfluxDb(mc.influxDb, device, data.UserEmail)
	fmt.Printf("Solar panel: name=%s, id=%d, changed switch to %t \n", device.Device.Name, device.Device.Id, device.IsOn)
}

func (mc *MQTTClient) HandleSPData(client mqtt.Client, msg mqtt.Message) {
	// Retrieve the last part of the split string, which is the device ID
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the JSON string into the struct
	valueStr := string(msg.Payload())
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return
	}
	device := mc.solarPanelRepository.Get(deviceId)
	mc.handleProduction(device, value)
	//mc.solarPanelRepository.UpdateSP(device)
	//mc.deviceRepository.Update(device.Device)
	fmt.Printf("Solar panel: name=%s, id=%d, generated value %f \n", device.Device.Name, device.Device.Id, value)
}

func (mc *MQTTClient) handleProduction(device models.SolarPanel, value float64) {
	batteries, err := mc.homeBatteryRepository.GetAllByEstateId(device.Device.RealEstate)
	if err != nil {
		return
	}
	surplus := 0.0
	if len(batteries) != 0 {
		valuePerBattery := value / float64(len(batteries))
		// value is divided between batteries and each battery takes the same value
		surplus = mc.calculateProductionForBatteries(batteries, device.Device.Id, valuePerBattery, false)
		if surplus != 0.0 {
			// surplus=what was left (if one of the batteries was full) is sent to batteries again (not divided)
			surplus = mc.calculateProductionForBatteries(batteries, device.Device.Id, surplus, true)
		}
	} else {
		saveSPDataToInfluxDb(mc.influxDb, device.Device.Id, "electrical_distribution", value)
	}
	if surplus != 0.0 {
		saveSPDataToInfluxDb(mc.influxDb, device.Device.Id, "electrical_distribution", surplus)
	}

}

func (mc *MQTTClient) calculateProductionForBatteries(batteries []models.HomeBattery, deviceId int, valuePerBattery float64, isSurplus bool) float64 {
	surplus := 0.0
	for i, _ := range batteries {
		if !batteries[i].Device.IsOnline {
			if !isSurplus {
				surplus = surplus + valuePerBattery
			}
			continue
		}
		if batteries[i].CurrentValue+valuePerBattery <= batteries[i].Size {
			batteries[i].CurrentValue = batteries[i].CurrentValue + valuePerBattery
			saveSPDataToInfluxDb(mc.influxDb, deviceId, strconv.Itoa(batteries[i].Device.Id), valuePerBattery)
			SaveHBDataToInfluxDb(mc.influxDb, batteries[i].Device.Id, batteries[i].CurrentValue)
			mc.homeBatteryRepository.Update(batteries[i])
			//everything that was supposed to go into the battery went into it
			if isSurplus {
				valuePerBattery = 0.0
				break
			}
		} else {
			//not everything that was supposed to go into the battery went into it
			produced := batteries[i].Size - batteries[i].CurrentValue
			batteries[i].CurrentValue = batteries[i].Size
			if isSurplus {
				valuePerBattery = valuePerBattery - produced
			} else {
				surplus = surplus + valuePerBattery - produced
			}
			saveSPDataToInfluxDb(mc.influxDb, deviceId, strconv.Itoa(batteries[i].Device.Id), produced)
			SaveHBDataToInfluxDb(mc.influxDb, batteries[i].Device.Id, batteries[i].CurrentValue)
			mc.homeBatteryRepository.Update(batteries[i])
		}
	}
	if isSurplus {
		return valuePerBattery
	} else {
		return surplus
	}
}

func saveSPSwitchChangeToInfluxDb(client influxdb2.Client, device models.SolarPanel, email string) {
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

	writeAPI.WritePoint(p)

	writeAPI.Flush()
	fmt.Println("Saved sp switch change to influxdb")
}

func saveSPDataToInfluxDb(client influxdb2.Client, deviceId int, batteryId string, value float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	p := influxdb2.NewPoint("solar_panel", //table
		map[string]string{"device_id": strconv.Itoa(deviceId), "battery_id": batteryId}, //tag
		map[string]interface{}{"electricity": value},                                    //field
		time.Now())

	// Write the point to InfluxDB
	writeAPI.WritePoint(p)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Saved sp data to influxdb")
}

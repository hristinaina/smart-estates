package mqtt_client

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	"strconv"
	"strings"
	"time"
)

func (mc *MQTTClient) HandleSprinklerMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("STIGLOOOOO")
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(msg.Topic())
	fmt.Println("DEVICE IDDDD " + strconv.Itoa(deviceId))
	_, err = mc.sprinkleRepository.UpdateIsOn(deviceId, true)
	if err != nil {
		fmt.Println(err)
	}
	saveSprinklerToInfluxDb(mc.influxDb, deviceId, "on", "auto")
}

func (mc *MQTTClient) HandleSprinklerOffMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("STIGLOOOOO2")
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("PRIMLJENA PORUKA2")
	fmt.Println(msg.Topic())
	fmt.Println("DEVICE IDDDD " + strconv.Itoa(deviceId))
	_, err = mc.sprinkleRepository.UpdateIsOn(deviceId, false)
	if err != nil {
		fmt.Println(err)
	}
	saveSprinklerToInfluxDb(mc.influxDb, deviceId, "off", "auto")
}

func saveSprinklerToInfluxDb(client influxdb2.Client, deviceId int, mode, user string) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	point := influxdb2.NewPoint("sprinkler", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{"action": mode, "user_id": user},
		time.Now()) // field

	writeAPI.WritePoint(point)
	writeAPI.Flush()
}

type SprinklerHistoryData struct {
	User   string
	Action string
}

func SprinklerQueryDeviceData(client influxdb2.Client, data dtos.ActionGraphRequest) map[string]SprinklerHistoryData {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := client.QueryAPI(Org)
	var query string

	if data.EndDate != "" && data.StartDate != "" && data.UserEmail == "none" {
		endDate, _ := time.Parse("2006-01-02", data.EndDate)
		endDate = endDate.AddDate(0, 0, 1)
		endDateStr := endDate.Format("2006-01-02")
		query = fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r._measurement == "sprinkler" and r.device_id == "%s")
    `, Bucket, data.StartDate, endDateStr, fmt.Sprint(data.DeviceId))
	} else if data.EndDate == "" && data.StartDate == "" && data.UserEmail != "none" {
		query = fmt.Sprintf(` 
        from(bucket: "%s")
        |> range(start: 0)
        |> filter(fn: (r) => r._measurement == "sprinkler" and r.device_id == "%s")
		|> filter(fn: (r) => r._field != "user_id" or (r._field == "user_id" and r._value == "%s"))
    `, Bucket, fmt.Sprint(data.DeviceId), data.UserEmail)
	} else if data.EndDate != "" && data.StartDate != "" && data.UserEmail != "none" {
		endDate, _ := time.Parse("2006-01-02", data.EndDate)
		endDate = endDate.AddDate(0, 0, 1)
		endDateStr := endDate.Format("2006-01-02")
		query = fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r._measurement == "sprinkler" and r.device_id == "%s")
		|> filter(fn: (r) => r._field != "user_id" or (r._field == "user_id" and r._value == "%s"))
    `, Bucket, data.StartDate, endDateStr, fmt.Sprint(data.DeviceId), data.UserEmail)
	} else {
		query = fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: 0)
        |> filter(fn: (r) => r._measurement == "sprinkler" and r.device_id == "%s")	
    `, Bucket, fmt.Sprint(data.DeviceId))
	}

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %s\n", err.Error())
		return nil
	}

	var resultPoints map[string]SprinklerHistoryData
	resultPoints = make(map[string]SprinklerHistoryData)
	localLocation, err := time.LoadLocation("Local")

	if err == nil {
		// Iterate over query response
		for result.Next() {
			localTime := result.Record().Time().In(localLocation)
			time := localTime.Format("2006-01-02 15:04:05")

			val, _ := resultPoints[time]

			switch field := result.Record().Field(); field {
			case "action":
				val.Action = result.Record().Value().(string)
			case "user_id":
				val.User = result.Record().Value().(string)
			default:
				fmt.Printf("unrecognized field %s.\n", field)
			}

			resultPoints[time] = val
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	for key, value := range resultPoints {
		if value.User == "" {
			delete(resultPoints, key)
		}
	}
	return resultPoints
}

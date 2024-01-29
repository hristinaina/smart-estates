package mqtt_client

import (
	"context"
	"encoding/json"
	"fmt"
	"smarthome-back/dtos"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (mc *MQTTClient) HandleWMAction(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	payload := string(msg.Payload())

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	mode := data["Mode"].(string)
	switchWM := data["Switch"].(bool)
	previous := data["Previous"].(string)
	user := data["UserEmail"].(string)
	fmt.Println("VES MASINA")
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(deviceId, data)

	saveWMToInfluxDb(mc.influxDb, deviceId, mode, previous, user, switchWM)
}

func saveWMToInfluxDb(client influxdb2.Client, deviceId int, mode, previous, user string, switchWM bool) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	action := "Turn off"

	if switchWM {
		action = "Turn on"
	}
	point := influxdb2.NewPoint("washing_machine", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{"action": action, "mode": mode, "user_id": user},
		time.Now()) // field

	writeAPI.WritePoint(point)
	writeAPI.Flush()

	time.Sleep(1 * time.Second)

	if previous != "" {
		point := influxdb2.NewPoint("washing_machine", // table
			map[string]string{"device_id": strconv.Itoa(deviceId)},                                        // tag
			map[string]interface{}{"action": "Turn off", "mode": previous, "user_id": "auto"}, time.Now()) // field

		writeAPI.WritePoint(point)
		writeAPI.Flush()
	}
	fmt.Println("Washing Machine influxdb")
}

type WMHistoryData struct {
	User   string
	Action string
	Mode   string
}

func GetWMHistory(client influxdb2.Client, data dtos.ActionGraphRequest) map[string]WMHistoryData {
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
        |> filter(fn: (r) => r._measurement == "washing_machine" and r.device_id == "%s")
    `, Bucket, data.StartDate, endDateStr, fmt.Sprint(data.DeviceId))
	} else if data.EndDate == "" && data.StartDate == "" && data.UserEmail != "none" {
		query = fmt.Sprintf(` 
        from(bucket: "%s")
        |> range(start: 0)
        |> filter(fn: (r) => r._measurement == "washing_machine" and r.device_id == "%s")
		|> filter(fn: (r) => r._field != "user_id" or (r._field == "user_id" and r._value == "%s"))
    `, Bucket, fmt.Sprint(data.DeviceId), data.UserEmail)
	} else if data.EndDate != "" && data.StartDate != "" && data.UserEmail != "none" {
		endDate, _ := time.Parse("2006-01-02", data.EndDate)
		endDate = endDate.AddDate(0, 0, 1)
		endDateStr := endDate.Format("2006-01-02")
		query = fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r._measurement == "washing_machine" and r.device_id == "%s")
		|> filter(fn: (r) => r._field != "user_id" or (r._field == "user_id" and r._value == "%s"))
    `, Bucket, data.StartDate, endDateStr, fmt.Sprint(data.DeviceId), data.UserEmail)
	} else {
		query = fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: 0)
        |> filter(fn: (r) => r._measurement == "washing_machine" and r.device_id == "%s")	
    `, Bucket, fmt.Sprint(data.DeviceId))
	}

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %s\n", err.Error())
		return nil
	}

	var resultPoints map[string]WMHistoryData
	resultPoints = make(map[string]WMHistoryData)
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
			case "mode":
				val.Mode = result.Record().Value().(string)
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

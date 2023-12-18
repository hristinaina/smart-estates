package mqtt_client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (mc *MQTTClient) HandleActionChange(client mqtt.Client, msg mqtt.Message) {
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
	switchAC := data["Switch"].(bool)
	temp := data["Temp"].(float64)
	previous := data["Previous"].(string)
	user := data["UserEmail"].(string)
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(deviceId, mode, temp, previous, user, switchAC)

	// todo sacuvaj u bazi
	saveACToInfluxDb(mc.influxDb, deviceId, mode, previous, user, switchAC)
	// fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
}

func saveACToInfluxDb(client influxdb2.Client, deviceId int, mode, previous, user string, switchAC bool) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	action := 0

	if switchAC {
		action = 1
	}

	point := influxdb2.NewPoint("air_conditioner1", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{"action": action, "mode": mode, "user_id": user},
		time.Now()) // field

	writeAPI.WritePoint(point)
	writeAPI.Flush()

	if previous != "" {
		point := influxdb2.NewPoint("air_conditioner1", // table
			map[string]string{"device_id": strconv.Itoa(deviceId)},                               // tag
			map[string]interface{}{"action": 0, "mode": previous, "user_id": "auto"}, time.Now()) // field
		fmt.Println(point)

		writeAPI.WritePoint(point)
		writeAPI.Flush()
	}
	fmt.Println("Air Conditioner influxdb")
}

type ACHistoryData struct {
	User   string
	Action int64
	Mode   string
}

func QueryDeviceData(client influxdb2.Client, deviceId string) map[time.Time]ACHistoryData {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := client.QueryAPI(Org)

	// Konstruisanje Flux upita
	query := fmt.Sprintf(`
        from(bucket: "%s")
        |> range(start: 0)
        |> filter(fn: (r) => r._measurement == "air_conditioner1" and r.device_id == "%s" )
    `, Bucket, deviceId)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %s\n", err.Error())
		return nil
	}

	var resultPoints map[time.Time]ACHistoryData
	resultPoints = make(map[time.Time]ACHistoryData)

	if err == nil {
		// Iterate over query response
		for result.Next() {

			val, _ := resultPoints[result.Record().Time()]
			// fmt.Println(val)

			switch field := result.Record().Field(); field {
			case "action":
				val.Action = result.Record().Value().(int64)
			case "mode":
				val.Mode = result.Record().Value().(string)
			case "user_id":
				val.User = result.Record().Value().(string)
				fmt.Println(val.User)
			default:
				fmt.Printf("unrecognized field %s.\n", field)
			}

			resultPoints[result.Record().Time()] = val

		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return resultPoints
}

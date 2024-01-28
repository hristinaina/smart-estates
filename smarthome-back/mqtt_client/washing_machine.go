package mqtt_client

import (
	"encoding/json"
	"fmt"
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
	temp := data["Temp"].(float64)
	previous := data["Previous"].(string)
	user := data["UserEmail"].(string)
	fmt.Println("PRIMLJENA PORUKA")
	fmt.Println(deviceId, mode, temp, previous, user, switchWM)

	saveWMToInfluxDb(mc.influxDb, deviceId, mode, previous, user, switchWM)
}

func saveWMToInfluxDb(client influxdb2.Client, deviceId int, mode, previous, user string, switchWM bool) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)
	action := 0

	if switchWM {
		action = 1
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
			map[string]string{"device_id": strconv.Itoa(deviceId)},                               // tag
			map[string]interface{}{"action": 0, "mode": previous, "user_id": "auto"}, time.Now()) // field

		writeAPI.WritePoint(point)
		writeAPI.Flush()
	}
	fmt.Println("Washing Machine influxdb")
}

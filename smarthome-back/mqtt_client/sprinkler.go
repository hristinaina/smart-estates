package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	fmt.Println("DEVICE IDDDD " + strconv.Itoa(deviceId))
	saveSprinklerToInfluxDb(mc.influxDb, deviceId, "on", "auto")
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

	fmt.Println("sprinkler influxdb")
}

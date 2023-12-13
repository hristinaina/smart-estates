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

var switchOn = false

var influxdb influxdb2.Client

// HandleHeartBeat callback function called when subscribed to TopicOnline. Update heartbeat time when "online" message is received
func (mc *MQTTClient) ReceiveValue(client mqtt.Client, msg mqtt.Message) {
	influxdb = mc.influxDb

	payload := string(msg.Payload())

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	deviceId := int(data["id"].(float64))
	temperature := data["temperature"].(float64)
	humidity := data["humidity"].(float64)

	saveValueToInfluxDb(mc.influxDb, deviceId, temperature, humidity)

	// GetLastOneHourValues(mc.influxDb, "7") // todo izbrisi ovo

	fmt.Printf("Ambient Sensor, id=%v, temeprature: %v °C, humidity: %v %% \n", deviceId, temperature, humidity)

}

func saveValueToInfluxDb(client influxdb2.Client, deviceId int, temperature, humidity float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	fmt.Println("PODACI")
	fmt.Println(deviceId)
	fmt.Println(temperature)
	fmt.Println(humidity)

	point := influxdb2.NewPoint("measurement1", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{"temperature": temperature, "humidity": humidity},
		time.Now()) // field

	// Write the point to InfluxDB
	writeAPI.WritePoint(point)

	// Close the write API to flush the buffer and release resources
	writeAPI.Flush()
	fmt.Println("Ambient sensor influxdb")
}

func (mc *MQTTClient) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	deviceId, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println(err)
	}
	status := string(msg.Payload())
	switchOn = status == "true"
	fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
}

type AmbientSensor struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func GetLastOneHourValues(deviceId string) map[time.Time]AmbientSensor {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := influxdb.QueryAPI(Org)

	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: -1h, stop: now())
	|> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")`,
		Bucket, deviceId)

	// result, err := queryAPI.QueryRaw(context.Background(), query, influxdb2.DefaultDialect())
	// if err == nil {
	// 	fmt.Println("QueryResult:")
	// 	fmt.Println(result)
	// } else {
	// 	panic(err)
	// }

	// return nil

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var resultPoints map[time.Time]AmbientSensor
	resultPoints = make(map[time.Time]AmbientSensor)

	// lista := []AmbientSensor{}

	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			// if result.TableChanged() {
			// 	fmt.Printf("table: %s\n", result.TableMetadata().String())
			// }

			val, ok := resultPoints[result.Record().Time()]

			if !ok {
				// val = models.Device{
				//     user: fmt.Sprintf("%v", result.Record().ValueByKey("user")),
				// }
			}

			switch field := result.Record().Field(); field {
			case "temperature":
				val.Temperature = result.Record().Value().(float64)
			case "humidity":
				val.Humidity = result.Record().Value().(float64)
			default:
				fmt.Printf("unrecognized field %s.\n", field)
			}

			// val.Timestamp = result.Record().Time()

			// lista = append(lista, val)

			resultPoints[result.Record().Time()] = val

		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	// fmt.Println("REZULTAT")
	// fmt.Println(resultPoints)
	// fmt.Println("LISTA")
	// fmt.Println(lista)

	return resultPoints
}

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

type AmbientSensor struct {
	Humidity    float64   `json:"humidity"`
	Temperature float64   `json:"temperature"`
	Timestemp   time.Time `json:"timestamp"`
}

var sensor AmbientSensor

func (mc *MQTTClient) ReceiveValue(client mqtt.Client, msg mqtt.Message) {
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

	setNewValue(temperature, humidity, time.Now())

	fmt.Printf("Ambient Sensor, id=%v, temeprature: %v °C, humidity: %v %% \n", deviceId, temperature, humidity)

}

func saveValueToInfluxDb(client influxdb2.Client, deviceId int, temperature, humidity float64) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

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

func GetLastOneHourValues(influxdb influxdb2.Client, deviceId string) map[time.Time]AmbientSensor {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := influxdb.QueryAPI(Org)

	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: -1h, stop: now())
	|> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")`,
		Bucket, deviceId)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var resultPoints map[time.Time]AmbientSensor
	resultPoints = make(map[time.Time]AmbientSensor)

	if err == nil {
		// Iterate over query response
		for result.Next() {

			val, _ := resultPoints[result.Record().Time()]

			switch field := result.Record().Field(); field {
			case "temperature":
				val.Temperature = result.Record().Value().(float64)
			case "humidity":
				val.Humidity = result.Record().Value().(float64)
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

	// fmt.Println("REZULTAT")
	// fmt.Println(resultPoints)
	// fmt.Println("LISTA")
	// fmt.Println(lista)

	return resultPoints
}

func setNewValue(temp, hmd float64, time time.Time) {
	sensor.Temperature = temp
	sensor.Humidity = hmd
	sensor.Timestemp = time
}

func GetNewValue() AmbientSensor {
	return sensor
}

func GetValuesForSelectedTime(influxdb influxdb2.Client, selectedTime, deviceId string) map[time.Time]AmbientSensor {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := influxdb.QueryAPI(Org)

	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: %s, stop: now())
	|> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")`,
		Bucket, selectedTime, deviceId)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var resultPoints map[time.Time]AmbientSensor
	resultPoints = make(map[time.Time]AmbientSensor)

	if err == nil {
		// Iterate over query response
		for result.Next() {

			val, _ := resultPoints[result.Record().Time()]

			switch field := result.Record().Field(); field {
			case "temperature":
				val.Temperature = result.Record().Value().(float64)
			case "humidity":
				val.Humidity = result.Record().Value().(float64)
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

	// fmt.Println("REZULTAT")
	// fmt.Println(resultPoints)
	// fmt.Println("LISTA")
	// fmt.Println(lista)

	return resultPoints
}

func GetValuesForDate(influxdb influxdb2.Client, start, end, deviceId string) map[time.Time]AmbientSensor {
	Org := "Smart Home"
	Bucket := "bucket"
	queryAPI := influxdb.QueryAPI(Org)

	fmt.Println("start")
	fmt.Println(start)
	fmt.Println("end")
	fmt.Println(end)

	query := fmt.Sprintf(`from(bucket:"%s") 
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")`,
		Bucket, start, end, deviceId)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var resultPoints map[time.Time]AmbientSensor
	resultPoints = make(map[time.Time]AmbientSensor)

	if err == nil {
		// Iterate over query response
		for result.Next() {

			val, _ := resultPoints[result.Record().Time()]

			switch field := result.Record().Field(); field {
			case "temperature":
				val.Temperature = result.Record().Value().(float64)
			case "humidity":
				val.Humidity = result.Record().Value().(float64)
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

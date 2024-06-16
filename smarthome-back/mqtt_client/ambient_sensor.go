package mqtt_client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// var switchOn = false

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

//	func (mc *MQTTClient) HandleSwitchChange(client mqtt.Client, msg mqtt.Message) {
//		parts := strings.Split(msg.Topic(), "/")
//		deviceId, err := strconv.Atoi(parts[len(parts)-1])
//		if err != nil {
//			fmt.Println(err)
//		}
//		status := string(msg.Payload())
//		switchOn = status == "true"
//		fmt.Printf("AmbientSensor id=%d, switch status: %s\n", deviceId, status)
//	}
func processingQuery(influxdb influxdb2.Client, query string) map[time.Time]AmbientSensor {
	// Initialize the InfluxDB query API
	queryAPI := influxdb.QueryAPI("Smart Home")

	// Execute the query
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	// Initialize a map to store the processed data
	resultPoints := make(map[time.Time]AmbientSensor)

	// Iterate over the query response
	for result.Next() {
		// Extract the record
		record := result.Record()

		// Retrieve or initialize the AmbientSensor value for the timestamp
		val, ok := resultPoints[record.Time()]
		if !ok {
			val = AmbientSensor{} // Initialize if not present
		}

		// Process each field from the query result
		switch field := record.Field(); field {
		case "temperature":
			if value, ok := record.Value().(float64); ok {
				val.Temperature = value
			} else {
				fmt.Printf("temperature field value is not a float64: %v\n", record.Value())
			}
		case "humidity":
			if value, ok := record.Value().(float64); ok {
				val.Humidity = value
			} else {
				fmt.Printf("humidity field value is not a float64: %v\n", record.Value())
			}
		default:
			fmt.Printf("unrecognized field %s.\n", field)
		}

		// Store the processed AmbientSensor value back into the map
		resultPoints[record.Time()] = val
	}

	// Check for errors during iteration
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	return resultPoints
}

func GetLastOneHourValues(influxdb influxdb2.Client, deviceId string) map[time.Time]AmbientSensor {
	query := fmt.Sprintf(`from(bucket:"bucket") 
		|> range(start: -1h, stop: now())
		|> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")
		|> aggregateWindow(every: 10m, fn: mean)`, deviceId)

	return processingQuery(influxdb, query)
}

func GetValuesForSelectedTime(influxdb influxdb2.Client, selectedTime, deviceId string) map[time.Time]AmbientSensor {
	query := fmt.Sprintf(`
        from(bucket:"bucket") 
        |> range(start: %s, stop: now())
        |> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")
        |> aggregateWindow(every: 1h, fn: mean)`, selectedTime, deviceId)

	return processingQuery(influxdb, query)
}

func GetValuesForDate(influxdb influxdb2.Client, start, end, deviceId string) map[time.Time]AmbientSensor {
	endDate, _ := time.Parse(time.RFC3339, end)
	endDate = endDate.AddDate(0, 0, 1)
	endDateStr := endDate.Format(time.RFC3339)

	query := fmt.Sprintf(`
        from(bucket:"bucket") 
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r._measurement == "measurement1" and r.device_id == "%s")
        |> aggregateWindow(every: 12h, fn: sum)`, start, endDateStr, deviceId)

	return processingQuery(influxdb, query)
}

func setNewValue(temp, hmd float64, time time.Time) {
	sensor.Temperature = temp
	sensor.Humidity = hmd
	sensor.Timestemp = time
}

func GetNewValue() AmbientSensor {
	return sensor
}

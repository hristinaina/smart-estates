package energetic

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	"smarthome-back/models/devices/energetic"
	repositories "smarthome-back/repositories/devices"
	"time"

	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type HomeBatteryService interface {
	Add(estate dtos.DeviceDTO) energetic.HomeBattery
	GetAllByEstateId(id int) ([]energetic.HomeBattery, error)
	Get(id int) energetic.HomeBattery
	GetConsumptionFromLastMinute(id int) (interface{}, error)
	GetConsumptionForLastHour(id int) interface{}
	GetConsumptionForSelectedTime(selectedTime string, estateId int) interface{}
	GetConsumptionForSelectedDate(startDate, endDate string, estateId int) interface{}
	GetStatusForSelectedTime(selectedTime string, estateId int) interface{}
	GetStatusForSelectedDate(startDate, endDate string, estateId int) interface{}
}

type HomeBatteryServiceImpl struct {
	db           *sql.DB
	repository   repositories.HomeBatteryRepository
	influxDb     influxdb2.Client
	cacheService cache.CacheService
}

func NewHomeBatteryService(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) HomeBatteryService {
	return &HomeBatteryServiceImpl{db: db, repository: repositories.NewHomeBatteryRepository(db, cacheService), influxDb: influxDb}
}

func (s *HomeBatteryServiceImpl) GetAllByEstateId(id int) ([]energetic.HomeBattery, error) {
	return s.repository.GetAllByEstateId(id)
}

func (s *HomeBatteryServiceImpl) Add(dto dtos.DeviceDTO) energetic.HomeBattery {
	return s.repository.Add(dto)
}

func (s *HomeBatteryServiceImpl) Get(id int) energetic.HomeBattery {
	return s.repository.Get(id)
}

func (s *HomeBatteryServiceImpl) GetConsumptionFromLastMinute(id int) (interface{}, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: -1m10s) 
			|> filter(fn: (r) => r["_measurement"] == "consumption"
			and r["_field"] == "electricity" and r["estate_id"] == "%d")
			|> yield(name: "sum")`, influxBucket, id)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return 0.0, err
	}

	var value float64
	// Iterate over query results
	for result.Next() {
		if result.Record().Value() != nil {
			value = value + result.Record().ValueByKey("_value").(float64)
		}
	}

	// Check for errors
	if result.Err() != nil {
		fmt.Printf("Error processing InfluxDB query results: %v\n", result.Err())
	}

	// Close the result set
	result.Close()
	return value, nil
}

func (s *HomeBatteryServiceImpl) GetConsumptionForSelectedTime(selectedTime string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: now())
	|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
 	|> aggregateWindow(every: 1h, fn: sum)`, selectedTime, estateId)

	return s.processingQuery(query)
}

func (s *HomeBatteryServiceImpl) GetConsumptionForSelectedDate(startDate, endDate string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "consumption" and r["_field"] == "electricity" and r["estate_id"] == "%d")
 	|> aggregateWindow(every: 12h, fn: sum)`, startDate, endDate, estateId)

	return s.processingQuery(query)
}

func (s *HomeBatteryServiceImpl) GetStatusForSelectedTime(selectedTime string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: now())
	|> filter(fn: (r) => r._measurement == "home_battery" and r["_field"] == "currentValue" and r["device_id"] == "%d")
	|> aggregateWindow(every: 1h, fn: mean)`, selectedTime, estateId)

	return s.processingQuery(query)
}

func (s *HomeBatteryServiceImpl) GetStatusForSelectedDate(startDate, endDate string, estateId int) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "home_battery" and r["_field"] == "currentValue" and r["device_id"] == "%d")
	|> aggregateWindow(every: 12h, fn: mean)`, startDate, endDate, estateId)

	return s.processingQuery(query)
}

func (s *HomeBatteryServiceImpl) GetConsumptionForLastHour(estateId int) interface{} {
	bucket := "bucket"

	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: -1h, stop: now()) 
			|> filter(fn: (r) => r["_measurement"] == "consumption"
			and r["_field"] == "electricity" and r["estate_id"] == "%d")
			|> aggregateWindow(every: 10m, fn: sum)`, bucket, estateId)
	return s.processingQuery(query)
}

func (s *HomeBatteryServiceImpl) processingQuery(query string) map[time.Time]float64 {
	Org := "Smart Home"
	queryAPI := s.influxDb.QueryAPI(Org)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing InfluxDB query:", err)
		return nil
	}

	var tempPoints map[string]float64
	tempPoints = make(map[string]float64)

	if err == nil {
		// Iterate over query response
		for result.Next() {
			if result.Record().Value() != nil {
				timeStr := result.Record().Time().Format("2006-01-02 15:04")
				tempPoints[timeStr] = tempPoints[timeStr] + result.Record().ValueByKey("_value").(float64)
			}
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	var resultPoints map[time.Time]float64
	resultPoints = make(map[time.Time]float64)

	for timeStr, value := range tempPoints {
		layout := "2006-01-02 15:04"

		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			fmt.Printf("Error parsing time '%s': %v\n", timeStr, err)
			continue
		}

		resultPoints[parsedTime] = value
	}

	return resultPoints
}

package services

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
	"smarthome-back/repositories"
)

type HomeBatteryService interface {
	Add(estate dto.DeviceDTO) models.HomeBattery
	GetAllByEstateId(id int) ([]models.HomeBattery, error)
	Get(id int) models.HomeBattery
	GetConsumptionFromLastMinute(id int) (interface{}, error)
}

type HomeBatteryServiceImpl struct {
	db         *sql.DB
	repository repositories.HomeBatteryRepository
	influxDb   influxdb2.Client
}

func NewHomeBatteryService(db *sql.DB, influxDb influxdb2.Client) HomeBatteryService {
	return &HomeBatteryServiceImpl{db: db, repository: repositories.NewHomeBatteryRepository(db), influxDb: influxDb}
}

func (s *HomeBatteryServiceImpl) GetAllByEstateId(id int) ([]models.HomeBattery, error) {
	return s.repository.GetAllByEstateId(id)
}

func (s *HomeBatteryServiceImpl) Add(dto dto.DeviceDTO) models.HomeBattery {
	return s.repository.Add(dto)
}

func (s *HomeBatteryServiceImpl) Get(id int) models.HomeBattery {
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

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
	"time"
)

type SolarPanelService interface {
	Add(estate dto.DeviceDTO) models.SolarPanel
	Get(id int) models.SolarPanel
	UpdateSP(device models.SolarPanel) bool
	GetGraphData(data dto.ActionGraphRequest) (dto.ActionGraphResponse, error)
	GetValueFromLastMinute(id int) (interface{}, error)
}

type SolarPanelServiceImpl struct {
	db         *sql.DB
	repository repositories.SolarPanelRepository
	influxDb   influxdb2.Client
}

func NewSolarPanelService(db *sql.DB, influxDb influxdb2.Client) SolarPanelService {
	return &SolarPanelServiceImpl{db: db, repository: repositories.NewSolarPanelRepository(db), influxDb: influxDb}
}

func (s *SolarPanelServiceImpl) Get(id int) models.SolarPanel {
	return s.repository.Get(id)
}

func (s *SolarPanelServiceImpl) Add(dto dto.DeviceDTO) models.SolarPanel {
	return s.repository.Add(dto)
}

func (s *SolarPanelServiceImpl) GetGraphData(data dto.ActionGraphRequest) (dto.ActionGraphResponse, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: %s, stop: %s) 
			|> filter(fn: (r) => r["_measurement"] == "solar_panel"
			and r["_field"] == "isOn" and r["user_id"] == "%s" and r["device_id"] == "%d")`, influxBucket,
		data.StartDate, data.EndDate, data.UserEmail, data.DeviceId)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return dto.ActionGraphResponse{}, err
	}

	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Error loading local time zone:", err)
		return dto.ActionGraphResponse{}, err
	}

	var response dto.ActionGraphResponse
	// Iterate over query results
	for result.Next() {
		if result.Record().Value() != nil {
			localTime := result.Record().Time().In(localLocation)
			response.Labels = append(response.Labels, localTime.Format("2006-01-02 15:04:05 MST"))
			response.Values = append(response.Values, result.Record().Value())
		}
	}

	// Check for errors
	if result.Err() != nil {
		fmt.Printf("Error processing InfluxDB query results: %v\n", result.Err())
	}

	// Close the result set
	result.Close()
	return response, nil
}

func (s *SolarPanelServiceImpl) UpdateSP(device models.SolarPanel) bool {
	return s.repository.UpdateSP(device)
}

func (s *SolarPanelServiceImpl) GetValueFromLastMinute(id int) (interface{}, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: -1m10s) 
			|> filter(fn: (r) => r["_measurement"] == "solar_panel"
			and r["_field"] == "electricity" and r["device_id"] == "%d")
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

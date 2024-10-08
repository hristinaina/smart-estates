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

type SolarPanelService interface {
	Add(estate dtos.DeviceDTO) energetic.SolarPanel
	Get(id int) energetic.SolarPanel
	UpdateSP(device energetic.SolarPanel) bool
	GetGraphData(data dtos.ActionGraphRequest) (dtos.ActionGraphResponse, error)
	GetValueFromLastMinute(id int) (interface{}, error)
	GetProductionForSP(request dtos.ActionGraphRequest) interface{}
}

type SolarPanelServiceImpl struct {
	db           *sql.DB
	repository   repositories.SolarPanelRepository
	influxDb     influxdb2.Client
	cacheService cache.CacheService
}

func NewSolarPanelService(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) SolarPanelService {
	return &SolarPanelServiceImpl{db: db, repository: repositories.NewSolarPanelRepository(db, cacheService), influxDb: influxDb}
}

func (s *SolarPanelServiceImpl) Get(id int) energetic.SolarPanel {
	return s.repository.Get(id)
}

func (s *SolarPanelServiceImpl) Add(dto dtos.DeviceDTO) energetic.SolarPanel {
	return s.repository.Add(dto)
}

func (s *SolarPanelServiceImpl) GetGraphData(data dtos.ActionGraphRequest) (dtos.ActionGraphResponse, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	query := ""
	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	if data.UserEmail == "all" {
		// Define your InfluxDB query with conditions
		query = fmt.Sprintf(`from(bucket:"%s")|> range(start: %s, stop: %s) 
			|> filter(fn: (r) => r["_measurement"] == "solar_panel"
			and r["_field"] == "isOn" and r["device_id"] == "%d")`, influxBucket,
			data.StartDate, data.EndDate, data.DeviceId)
	} else {
		query = fmt.Sprintf(`from(bucket:"%s")|> range(start: %s, stop: %s) 
			|> filter(fn: (r) => r["_measurement"] == "solar_panel"
			and r["_field"] == "isOn" and r["user_id"] == "%s" and r["device_id"] == "%d")`, influxBucket,
			data.StartDate, data.EndDate, data.UserEmail, data.DeviceId)
	}

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return dtos.ActionGraphResponse{}, err
	}

	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Error loading local time zone:", err)
		return dtos.ActionGraphResponse{}, err
	}

	var response dtos.ActionGraphResponse
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

func (s *SolarPanelServiceImpl) UpdateSP(device energetic.SolarPanel) bool {
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

func (s *SolarPanelServiceImpl) GetProductionForSP(data dtos.ActionGraphRequest) interface{} {
	query := fmt.Sprintf(`from(bucket:"bucket") 
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "solar_panel" and r["_field"] == "electricity" and r["device_id"] == "%d")
		|> aggregateWindow(every: 12h, fn: sum)`, data.StartDate, data.EndDate, data.DeviceId)
	return s.processingQuery(query)
}

func (s *SolarPanelServiceImpl) processingQuery(query string) map[time.Time]float64 {
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

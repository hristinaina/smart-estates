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

type SolarPanelService interface {
	Add(estate dto.DeviceDTO) models.SolarPanel
	Get(id int) models.SolarPanel
	UpdateSP(device models.SolarPanel) bool
	GetGraphData(data dto.ActionGraph) (interface{}, error)
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

func (s *SolarPanelServiceImpl) GetGraphData(data dto.ActionGraph) (interface{}, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	fmt.Println(data)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: %d, stop: %d) 
			|> filter(fn: (r) => r["_measurement"] == "solar_panel"
			and r["_field"] == "isOn" and r["user_id"] == "%s" and r["device_id"] == "%d")`, influxBucket,
		data.StartDate, data.EndDate, data.UserEmail, data.DeviceId)
	// todo konvertovati datume ??

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return nil, err
	}

	// Iterate over query results
	for result.Next() {
		// Process the result, e.g., print data points
		fmt.Printf("Record: %v\n", result.Record().Values)
	}

	// Check for errors
	if result.Err() != nil {
		fmt.Printf("Error processing InfluxDB query results: %v\n", result.Err())
	}

	// Close the result set
	result.Close()
	return nil, nil
}

func (s *SolarPanelServiceImpl) UpdateSP(device models.SolarPanel) bool {
	return s.repository.UpdateSP(device)
}

package energetic

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	"smarthome-back/models/devices/energetic"
	repositories "smarthome-back/repositories/devices"
	"strconv"
	"time"
)

type EVChargerService interface {
	Add(estate dtos.DeviceDTO) energetic.EVCharger
	Get(id int) energetic.EVCharger
	GetLastPercentage(id int) (float64, error)
	GetHistoryActions(data dtos.ActionGraphRequest) (map[time.Time]EVActionHistory, error)
}

type EVChargerServiceImpl struct {
	db         *sql.DB
	influxDb   influxdb2.Client
	repository repositories.EVChargerRepository
}

func NewEVChargerService(db *sql.DB, influxdb influxdb2.Client) EVChargerService {
	return &EVChargerServiceImpl{db: db, influxDb: influxdb, repository: repositories.NewEVChargerRepository(db)}
}

func (s *EVChargerServiceImpl) Get(id int) energetic.EVCharger {
	return s.repository.Get(id)
}

func (s *EVChargerServiceImpl) Add(dto dtos.DeviceDTO) energetic.EVCharger {
	return s.repository.Add(dto)
}

func (s *EVChargerServiceImpl) GetLastPercentage(id int) (float64, error) {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := s.influxDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -1h)
			|> filter(fn: (r) =>
				r._measurement == "ev_charger" and
				r.device_id == "%s" and
				r.action == "%s"
			)
			|> last()
	`, influxBucket, strconv.Itoa(id), "percentageChange")

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 90, fmt.Errorf("error executing query: %v", err)
	}
	defer result.Close()

	if result.Next() {
		point := result.Record().Values()
		value, ok := point["_value"].(float64)
		if !ok {
			return 90, fmt.Errorf("unable to extract value from the result")
		}
		return value, nil
	}

	return 90, fmt.Errorf("no data found for device_id %s and action %s", strconv.Itoa(id), "percentageChange")
}

type EVActionHistory struct {
	User       string
	Action     string
	Percentage float64
	Plug       int
}

func (s *EVChargerServiceImpl) GetHistoryActions(data dtos.ActionGraphRequest) (map[time.Time]EVActionHistory, error) {
	Org := "Smart Home"
	Bucket := "bucket"

	queryAPI := s.influxDb.QueryAPI(Org)
	query := ""
	if data.UserEmail != "all" {
		query = fmt.Sprintf(`
	from(bucket:"%s")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "ev_charger" and r["device_id"] == "%s" and r["user_id"] == "%s")
	|> sort(columns: ["_time"])
	`, Bucket, data.StartDate, data.EndDate, strconv.Itoa(data.DeviceId), data.UserEmail)
	} else {
		query = fmt.Sprintf(`
	from(bucket:"%s")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "ev_charger" and r["device_id"] == "%s")
	|> sort(columns: ["_time"])
	`, Bucket, data.StartDate, data.EndDate, strconv.Itoa(data.DeviceId))
	}
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error executing InfluxDB query: %v", err)
	}
	defer result.Close()

	// Map to store the result
	actionHistoryMap := make(map[time.Time]EVActionHistory)

	// Iterate over result records
	for result.Next() {
		record := result.Record()
		timestampUTC := result.Record().Time().UnixNano()
		timestamp := time.Unix(0, timestampUTC).UTC()
		values := record.Values()

		// Extract relevant fields
		user := values["user_id"].(string)
		action := values["action"].(string)
		percentage := values["_value"].(float64)
		plug := values["plug_id"].(string)

		plugId, err := strconv.Atoi(plug)
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error:", err)
		}

		// Create EVActionHistory struct
		actionHistory := EVActionHistory{
			User:       user,
			Action:     action,
			Percentage: percentage,
			Plug:       plugId,
		}

		// Store in the map
		actionHistoryMap[timestamp] = actionHistory
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("error retrieving query result: %v", result.Err())
	}

	return actionHistoryMap, nil
}

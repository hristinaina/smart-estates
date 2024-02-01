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
)

type EVChargerService interface {
	Add(estate dtos.DeviceDTO) energetic.EVCharger
	Get(id int) energetic.EVCharger
	GetLastPercentage(id int) (float64, error)
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

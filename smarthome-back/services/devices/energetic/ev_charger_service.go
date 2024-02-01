package energetic

import (
	"database/sql"
	_ "database/sql"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	"smarthome-back/models/devices/energetic"
	repositories "smarthome-back/repositories/devices"
)

type EVChargerService interface {
	Add(estate dtos.DeviceDTO) energetic.EVCharger
	Get(id int) energetic.EVCharger
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

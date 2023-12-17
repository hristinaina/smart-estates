package services

import (
	"database/sql"
	_ "database/sql"
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

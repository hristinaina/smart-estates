package services

import (
	"database/sql"
	_ "database/sql"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
	"smarthome-back/repositories"
)

type HomeBatteryService interface {
	Add(estate dto.DeviceDTO) models.HomeBattery
	GetAllByEstateId(id int) ([]models.HomeBattery, error)
}

type HomeBatteryServiceImpl struct {
	db         *sql.DB
	repository repositories.HomeBatteryRepository
}

func NewHomeBatteryService(db *sql.DB) HomeBatteryService {
	return &HomeBatteryServiceImpl{db: db, repository: repositories.NewHomeBatteryRepository(db)}
}

func (s *HomeBatteryServiceImpl) GetAllByEstateId(id int) ([]models.HomeBattery, error) {
	return s.repository.GetAllByEstateId(id)
}

func (s *HomeBatteryServiceImpl) Add(dto dto.DeviceDTO) models.HomeBattery {
	return s.repository.Add(dto)
}

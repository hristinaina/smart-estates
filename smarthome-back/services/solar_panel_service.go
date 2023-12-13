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

type SolarPanelService interface {
	Add(estate dto.DeviceDTO) models.SolarPanel
	Get(id int) models.SolarPanel
	UpdateSP(device models.SolarPanel) bool
}

type SolarPanelServiceImpl struct {
	db         *sql.DB
	repository repositories.SolarPanelRepository
}

func NewSolarPanelService(db *sql.DB) SolarPanelService {
	return &SolarPanelServiceImpl{db: db, repository: repositories.NewSolarPanelRepository(db)}
}

func (s *SolarPanelServiceImpl) Get(id int) models.SolarPanel {
	return s.repository.Get(id)
}

func (s *SolarPanelServiceImpl) Add(dto dto.DeviceDTO) models.SolarPanel {
	return s.repository.Add(dto)
}

func (s *SolarPanelServiceImpl) UpdateSP(device models.SolarPanel) bool {
	return s.repository.UpdateSP(device)
}

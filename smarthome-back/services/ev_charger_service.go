package services

import (
	"database/sql"
	_ "database/sql"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
)

type EVChargerService interface {
	Add(estate dto.DeviceDTO) models.EVCharger
}

type EVChargerServiceImpl struct {
	db *sql.DB
}

func NewEVChargerService(db *sql.DB) EVChargerService {
	return &EVChargerServiceImpl{db: db}
}

func (s *EVChargerServiceImpl) Add(dto dto.DeviceDTO) models.EVCharger {
	// TODO: add some validation and exception throwing
	device := dto.ToEVCharger()
	tx, err := s.db.Begin()
	if err != nil {
		return models.EVCharger{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return models.EVCharger{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.EVCharger{}
	}

	// Insert the new EVCharger into the EVCharger table
	result, err = tx.Exec(`
		INSERT INTO EvCharger (DeviceId, ChargingPower, Connections)
		VALUES (?, ?, ?)
	`, deviceID, device.ChargingPower, device.Connections)
	if err != nil {
		return models.EVCharger{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.EVCharger{}
	}
	device.Device.Id = int(deviceID)
	return device
}

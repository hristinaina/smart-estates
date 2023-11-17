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

type AirConditionerService interface {
	Add(estate dto.DeviceDTO) models.AirConditioner
}

type AirConditionerServiceImpl struct {
	db *sql.DB
}

func NewAirConditionerService(db *sql.DB) AirConditionerService {
	return &AirConditionerServiceImpl{db: db}
}

func (s *AirConditionerServiceImpl) Add(dto dto.DeviceDTO) models.AirConditioner {
	// TODO: add some validation and exception throwing
	device := dto.ToAirConditioner()
	tx, err := s.db.Begin()
	if err != nil {
		return models.AirConditioner{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, Picture, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?, ?)
	`, device.Device.Device.Name, device.Device.Device.Type, device.Device.Device.Picture, device.Device.Device.RealEstate,
		device.Device.Device.IsOnline)
	if err != nil {
		return models.AirConditioner{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.AirConditioner{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.Device.PowerSupply, device.Device.PowerConsumption)
	if err != nil {
		return models.AirConditioner{}
	}

	// Insert the new air conditioner into the AirConditioner table
	result, err = tx.Exec(`
		INSERT INTO AirConditioner (DeviceId, MinTemperature, MaxTemperature)
		VALUES (?, ?, ?)
	`, deviceID, device.MinTemperature, device.MaxTemperature)
	if err != nil {
		return models.AirConditioner{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.AirConditioner{}
	}
	device.Device.Device.Id = int(deviceID)
	return device
}

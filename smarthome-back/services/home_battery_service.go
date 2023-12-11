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

type HomeBatteryService interface {
	Add(estate dto.DeviceDTO) models.HomeBattery
}

type HomeBatteryServiceImpl struct {
	db *sql.DB
}

func NewHomeBatteryService(db *sql.DB) HomeBatteryService {
	return &HomeBatteryServiceImpl{db: db}
}

func (s *HomeBatteryServiceImpl) Add(dto dto.DeviceDTO) models.HomeBattery {
	// TODO: add some validation and exception throwing
	device := dto.ToHomeBattery()
	tx, err := s.db.Begin()
	if err != nil {
		return models.HomeBattery{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return models.HomeBattery{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.HomeBattery{}
	}

	// Insert the new Home Battery into the Home Battery table
	result, err = tx.Exec(`
		INSERT INTO HomeBattery (DeviceId, Size)
		VALUES (?, ?)
	`, deviceID, device.Size)
	if err != nil {
		return models.HomeBattery{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.HomeBattery{}
	}
	device.Device.Id = int(deviceID)
	return device
}

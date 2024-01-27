package energetic

import (
	"database/sql"
	_ "database/sql"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/dtos"
	"smarthome-back/models/devices/energetic"
)

type EVChargerService interface {
	Add(estate dtos.DeviceDTO) energetic.EVCharger
}

type EVChargerServiceImpl struct {
	db *sql.DB
}

func NewEVChargerService(db *sql.DB) EVChargerService {
	return &EVChargerServiceImpl{db: db}
}

func (s *EVChargerServiceImpl) Add(dto dtos.DeviceDTO) energetic.EVCharger {
	// TODO: add some validation and exception throwing
	device := dto.ToEVCharger()
	tx, err := s.db.Begin()
	if err != nil {
		return energetic.EVCharger{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return energetic.EVCharger{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return energetic.EVCharger{}
	}

	// Insert the new EVCharger into the EVCharger table
	result, err = tx.Exec(`
		INSERT INTO EvCharger (DeviceId, ChargingPower, Connections)
		VALUES (?, ?, ?)
	`, deviceID, device.ChargingPower, device.Connections)
	if err != nil {
		return energetic.EVCharger{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return energetic.EVCharger{}
	}
	device.Device.Id = int(deviceID)
	return device
}

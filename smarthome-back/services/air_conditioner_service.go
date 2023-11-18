package services

import (
	"database/sql"
	"fmt"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
)

type AirConditionerService interface {
	Add(estate dto.DeviceDTO) models.AirConditioner
	Get(id int) models.AirConditioner
}

type AirConditionerServiceImpl struct {
	db *sql.DB
}

func NewAirConditionerService(db *sql.DB) AirConditionerService {
	return &AirConditionerServiceImpl{db: db}
}

func (s *AirConditionerServiceImpl) Get(id int) models.AirConditioner {
	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.Picture,
			Device.RealEstate,
			Device.IsOnline,
			ConsumptionDevice.PowerSupply,
			ConsumptionDevice.PowerConsumption,
			AirConditioner.MinTemperature,
			AirConditioner.MaxTemperature
		FROM
			AirConditioner
		JOIN ConsumptionDevice ON AirConditioner.DeviceId = ConsumptionDevice.DeviceId
		JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
		WHERE
			Device.Id = ?
	`

	// Execute the query
	row := s.db.QueryRow(query, id)

	var ac models.AirConditioner
	var device models.Device
	var consDevice models.ConsumptionDevice

	err := row.Scan(
		&device.Id,
		&device.Name,
		&device.Type,
		&device.Picture,
		&device.RealEstate,
		&device.IsOnline,
		&consDevice.PowerSupply,
		&consDevice.PowerConsumption,
		&ac.MinTemperature,
		&ac.MaxTemperature,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No air conditioner found with the specified ID")
		} else {
			fmt.Println("Error retrieving air conditioner:", err)
		}
		return models.AirConditioner{}
	}
	consDevice.Device = device
	ac.Device = consDevice
	return ac
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

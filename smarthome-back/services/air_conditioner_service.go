package services

import (
	"database/sql"
	"fmt"
	"smarthome-back/dto"
	models "smarthome-back/models/devices"
	"time"
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

// func (s *AirConditionerServiceImpl) Get(id int) models.AirConditioner {
// 	// todo change this
// 	query := `
// 		SELECT
// 			Device.Id,
// 			Device.Name,
// 			Device.Type,
// 			Device.RealEstate,
// 			Device.IsOnline,
// 			Device.StatusTimeStamp,
// 			ConsumptionDevice.PowerSupply,
// 			ConsumptionDevice.PowerConsumption,
// 			AirConditioner.MinTemperature,
// 			AirConditioner.MaxTemperature
// 		FROM
// 			AirConditioner
// 		JOIN ConsumptionDevice ON AirConditioner.DeviceId = ConsumptionDevice.DeviceId
// 		JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
// 		WHERE
// 			Device.Id = ?
// 	`

// 	// Execute the query
// 	row := s.db.QueryRow(query, id)

// 	var ac models.AirConditioner
// 	var device models.Device
// 	var consDevice models.ConsumptionDevice

//		err := row.Scan(
//			&device.Id,
//			&device.Name,
//			&device.Type,
//			&device.RealEstate,
//			&device.IsOnline,
//			&device.StatusTimeStamp,
//			&consDevice.PowerSupply,
//			&consDevice.PowerConsumption,
//			&ac.MinTemperature,
//			&ac.MaxTemperature,
//		)
//		if err != nil {
//			if err == sql.ErrNoRows {
//				fmt.Println("No air conditioner found with the specified ID")
//			} else {
//				fmt.Println("Error retrieving air conditioner:", err)
//			}
//			return models.AirConditioner{}
//		}
//		consDevice.Device = device
//		ac.Device = consDevice
//		return ac
//	}
func (s *AirConditionerServiceImpl) Get(id int) models.AirConditioner {
	fmt.Println("USLOOOOOOOOOOO")
	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			ConsumptionDevice.PowerSupply,
			ConsumptionDevice.PowerConsumption,
			AirConditioner.MinTemperature,
			AirConditioner.MaxTemperature,
			AirConditioner.Mode,
			SpecialModes.StartTime,
			SpecialModes.EndTime,
			SpecialModes.Temperature,
			SpecialModes.SelectedDays
		FROM
			AirConditioner
		JOIN ConsumptionDevice ON AirConditioner.DeviceId = ConsumptionDevice.DeviceId
		JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
		LEFT JOIN SpecialModes ON AirConditioner.DeviceId = SpecialModes.DeviceId
		WHERE
			Device.Id = ?
	`

	// Execute the query
	rows, err := s.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return models.AirConditioner{}
	}
	defer rows.Close()

	var ac models.AirConditioner
	var device models.Device
	var consDevice models.ConsumptionDevice
	var specialModes []models.SpecialMode

	for rows.Next() {
		var startTimeStr, endTimeStr string
		var selectedDays string
		var temperature float32

		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Type,
			&device.RealEstate,
			&device.IsOnline,
			&device.StatusTimeStamp,
			&consDevice.PowerSupply,
			&consDevice.PowerConsumption,
			&ac.MinTemperature,
			&ac.MaxTemperature,
			&ac.Mode,
			&startTimeStr,
			&endTimeStr,
			&temperature,
			&selectedDays,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return models.AirConditioner{}
		}

		startTime, err := time.Parse("15:04:05", startTimeStr)
		if err != nil {
			fmt.Println("Error parsing StartTime:", err)
			return models.AirConditioner{}
		}

		endTime, err := time.Parse("15:04:05", endTimeStr)
		if err != nil {
			fmt.Println("Error parsing StartTime:", err)
			return models.AirConditioner{}
		}

		consDevice.Device = device
		ac.Device = consDevice

		// Dodajte svaki red rezultata kao poseban SpecialMode
		specialMode := models.SpecialMode{
			StartTime:    startTime,
			EndTime:      endTime,
			Temperature:  temperature,
			SelectedDays: selectedDays,
		}
		specialModes = append(specialModes, specialMode)
	}

	ac.SpecialMode = specialModes

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
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Device.Name, device.Device.Device.Type, device.Device.Device.RealEstate,
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

	// todo add mode
	// Insert the new air conditioner into the AirConditioner table
	result, err = tx.Exec(`
		INSERT INTO AirConditioner (DeviceId, MinTemperature, MaxTemperature)
		VALUES (?, ?, ?)
	`, deviceID, device.MinTemperature, device.MaxTemperature)
	if err != nil {
		return models.AirConditioner{}
	}

	// todo add special mode

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.AirConditioner{}
	}
	device.Device.Device.Id = int(deviceID)
	return device
}

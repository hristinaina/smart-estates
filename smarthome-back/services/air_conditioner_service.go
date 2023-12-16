package services

import (
	"database/sql"
	"fmt"
	"smarthome-back/dto"
	models "smarthome-back/models/devices"
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
			SpecialModes.Mode,
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
		var startTimeStr, endTimeStr sql.NullString
		var mode sql.NullString
		var selectedDays sql.NullString
		var temperature sql.NullFloat64

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
			&mode,
			&temperature,
			&selectedDays,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return models.AirConditioner{}
		}

		consDevice.Device = device
		ac.Device = consDevice

		specialMode := models.SpecialMode{}
		if startTimeStr.Valid {
			specialMode.StartTime = startTimeStr.String
		}
		if endTimeStr.Valid {
			specialMode.EndTime = endTimeStr.String
		}
		if mode.Valid {
			specialMode.Mode = mode.String
		}
		if temperature.Valid {
			specialMode.Temperature = float32(temperature.Float64)
		}
		if selectedDays.Valid {
			specialMode.SelectedDays = selectedDays.String
		}
		specialModes = append(specialModes, specialMode)
	}

	ac.SpecialMode = specialModes

	return ac
}

func (s *AirConditionerServiceImpl) Add(dto dto.DeviceDTO) models.AirConditioner {
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
		fmt.Println(err)
		return models.AirConditioner{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return models.AirConditioner{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.Device.PowerSupply, device.Device.PowerConsumption)
	if err != nil {
		fmt.Println(err)
		return models.AirConditioner{}
	}

	// Insert the new air conditioner into the AirConditioner table
	result, err = tx.Exec(`
		INSERT INTO AirConditioner (DeviceId, MinTemperature, MaxTemperature, Mode)
		VALUES (?, ?, ?, ?)
	`, deviceID, device.MinTemperature, device.MaxTemperature, device.Mode)
	if err != nil {
		fmt.Println("ovde je greska")
		fmt.Println(err)
		return models.AirConditioner{}
	}

	if len(device.SpecialMode) != 0 {
		for _, mode := range device.SpecialMode {
			result, err = tx.Exec(`
			INSERT INTO specialModes (DeviceId, StartTime, EndTime, Mode, Temperature, SelectedDays)
			VALUES (?, ?, ?, ?, ?, ?)
		`, deviceID, mode.StartTime, mode.EndTime, mode.Mode, mode.Temperature, mode.SelectedDays)
			if err != nil {
				fmt.Println("ovde jeeeeeeeeeee")
				fmt.Println(err)
				return models.AirConditioner{}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return models.AirConditioner{}
	}

	device.Device.Device.Id = int(deviceID)
	return device
}

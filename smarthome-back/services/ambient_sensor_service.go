package services

import (
	"database/sql"
	"fmt"
	"smarthome-back/dto"
	models "smarthome-back/models/devices"
)

type AmbientSensorService interface {
	Get(id int) (models.ConsumptionDevice, error)
	Add(dto dto.DeviceDTO) models.ConsumptionDevice
}

type AmbientSensorServiceImpl struct {
	db *sql.DB
}

func NewAmbientSensorService(db *sql.DB) AmbientSensorService {
	return &AmbientSensorServiceImpl{db: db}
}

func (as *AmbientSensorServiceImpl) Get(id int) (models.ConsumptionDevice, error) {
	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			ConsumptionDevice.PowerSupply,
			ConsumptionDevice.PowerConsumption
		FROM
			ConsumptionDevice
		JOIN 
			Device ON ConsumptionDevice.DeviceId = Device.Id
		WHERE
			Device.Id = ?;
	`
	// Execute the query
	rows, err := as.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return models.ConsumptionDevice{}, err
	}
	defer rows.Close()

	var device models.Device
	var consDevice models.ConsumptionDevice

	for rows.Next() {
		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Type,
			&device.RealEstate,
			&device.IsOnline,
			&device.StatusTimeStamp,
			&consDevice.PowerSupply,
			&consDevice.PowerConsumption,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return models.ConsumptionDevice{}, err
		}
		consDevice.Device = device
	}
	return consDevice, nil
}

func (as *AmbientSensorServiceImpl) Add(dto dto.DeviceDTO) models.ConsumptionDevice {
	device := dto.ToAmbientSensor()
	tx, err := as.db.Begin()
	if err != nil {
		return models.ConsumptionDevice{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.PowerSupply, device.PowerConsumption)
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	device.Device.Id = int(deviceID)
	return device
}

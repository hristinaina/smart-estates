package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"smarthome-back/dto"
	models "smarthome-back/models/devices"
)

type HomeBatteryRepository interface {
	Add(estate dto.DeviceDTO) models.HomeBattery
	GetAllByEstateId(id int) ([]models.HomeBattery, error)
	Update(device models.HomeBattery) bool
}

type HomeBatteryRepositoryImpl struct {
	db *sql.DB
}

func NewHomeBatteryRepository(db *sql.DB) HomeBatteryRepository {
	return &HomeBatteryRepositoryImpl{db: db}
}

func (s *HomeBatteryRepositoryImpl) GetAllByEstateId(id int) ([]models.HomeBattery, error) {
	query := `
		SELECT
			d.id,
			d.name,
			d.realEstate,
			d.isOnline,
			hb.size,
			hb.currentValue
		FROM
			device d
		JOIN
			homeBattery hb ON d.id = hb.deviceId
		WHERE
			d.realEstate = ?
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set
	var batteries []models.HomeBattery
	for rows.Next() {
		var device models.Device
		var hb models.HomeBattery

		//todo da li treba da scan bude skroz ispunjen?
		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.RealEstate,
			&device.IsOnline,
			&hb.Size,
			&hb.CurrentValue,
		)
		if err != nil {
			log.Fatal(err)
		}

		hb.Device = device
		batteries = append(batteries, hb)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return batteries, nil
}

func (s *HomeBatteryRepositoryImpl) Add(dto dto.DeviceDTO) models.HomeBattery {
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

func (res *HomeBatteryRepositoryImpl) Update(device models.HomeBattery) bool {
	query := "UPDATE homeBattery SET currentValue = ? WHERE deviceId = ?"
	_, err := res.db.Exec(query, device.CurrentValue, device.Device.Id)
	if err != nil {
		fmt.Println("Failed to update device:", err)
		return false
	}
	return true
}

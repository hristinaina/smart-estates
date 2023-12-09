package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/models/devices"
)

type DeviceRepository interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	GetAll() []models.Device
	GetDevicesByUserID(userID int) ([]models.Device, error)
	Update(device models.Device) bool
}

type DeviceRepositoryImpl struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) DeviceRepository {
	return &DeviceRepositoryImpl{db: db}
}

func (res *DeviceRepositoryImpl) GetAll() []models.Device {
	query := "SELECT * FROM device"
	rows, err := res.db.Query(query)
	if CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate,
			&device.IsOnline, &device.StatusTimeStamp); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Device{}
		}
		devices = append(devices, device)
	}

	return devices
}

func (res *DeviceRepositoryImpl) GetAllByEstateId(estateId int) []models.Device {
	query := "SELECT * FROM device WHERE REALESTATE = ?"
	rows, err := res.db.Query(query, estateId)
	if CheckIfError(err) {
		//todo raise an exception and catch it in controller?
		return nil
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		if err := rows.Scan(&device.Id, &device.Name, &device.Type,
			&device.RealEstate, &device.IsOnline, &device.StatusTimeStamp); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Device{}
		}
		devices = append(devices, device)
	}

	return devices
}

func (res *DeviceRepositoryImpl) Get(id int) (models.Device, error) {
	query := "SELECT * FROM device WHERE ID = ?"
	rows, err := res.db.Query(query, id)

	if CheckIfError(err) {
		return models.Device{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var (
			device models.Device
		)
		if err := rows.Scan(&device.Id, &device.Name, &device.Type,
			&device.RealEstate, &device.IsOnline, &device.StatusTimeStamp); err != nil {
			fmt.Println("Error: ", err.Error())
			return models.Device{}, err
		}
		return device, nil
	}
	return models.Device{}, err
}

func (res *DeviceRepositoryImpl) GetDevicesByUserID(userID int) ([]models.Device, error) {
	// Perform a database query to get devices by user ID
	rows, err := res.db.Query("SELECT id, name, type, realestate, isonline FROM device WHERE realestate IN (SELECT id FROM realestate WHERE userid = ?)", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the query result and populate the devices slice
	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate, &device.IsOnline)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (res *DeviceRepositoryImpl) Update(device models.Device) bool {
	query := "UPDATE device SET name = ?, type = ?, realestate = ?, isonline = ?, statustimestamp = ? WHERE id = ?"
	_, err := res.db.Exec(query, device.Name, device.Type, device.RealEstate, device.IsOnline, device.StatusTimeStamp, device.Id)
	if err != nil {
		fmt.Println("Failed to update device:", err)
		return false
	}
	return true
}

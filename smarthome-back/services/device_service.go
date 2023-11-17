package services

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/models/devices"
)

type DeviceService interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
}

type DeviceServiceImpl struct {
	db *sql.DB
}

func NewDeviceService(db *sql.DB) DeviceService {
	return &DeviceServiceImpl{db: db}
}

func (res *DeviceServiceImpl) GetAllByEstateId(estateId int) []models.Device {
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
			&device.Picture, &device.RealEstate, &device.IsOnline); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Device{}
			//todo raise an exception and catch it in controller?
		}
		devices = append(devices, device)
		fmt.Println(device)
	}

	return devices
}

func (res *DeviceServiceImpl) Get(id int) (models.Device, error) {
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
			&device.Picture, &device.RealEstate, &device.IsOnline); err != nil {
			fmt.Println("Error: ", err.Error())
			return models.Device{}, err
		}
		return device, nil
	}
	return models.Device{}, err
}

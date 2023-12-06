package services

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
	"smarthome-back/mqtt_client"
	"strconv"
)

type DeviceService interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	Add(estate dto.DeviceDTO) (models.Device, error)
	GetAll() []models.Device
}

type DeviceServiceImpl struct {
	db                    *sql.DB
	airConditionerService AirConditionerService
	evChargerService      EVChargerService
	homeBatteryService    HomeBatteryService
	mqtt                  *mqtt_client.MQTTClient
}

// todo send mqtt to all device_services
func NewDeviceService(db *sql.DB, mqtt *mqtt_client.MQTTClient) DeviceService {
	return &DeviceServiceImpl{db: db, airConditionerService: NewAirConditionerService(db), evChargerService: NewEVChargerService(db),
		homeBatteryService: NewHomeBatteryService(db), mqtt: mqtt}
}

func (res *DeviceServiceImpl) GetAll() []models.Device {
	query := "SELECT * FROM device"
	rows, err := res.db.Query(query)
	if CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.Picture, &device.RealEstate,
			&device.IsOnline); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Device{}
		}
		devices = append(devices, device)
		fmt.Println(device)
	}

	return devices
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

func (res *DeviceServiceImpl) Add(dto dto.DeviceDTO) (models.Device, error) {
	devices, err := res.getDevicesByUserID(dto.UserId)
	if err != nil {
		return models.Device{}, err
	}
	for _, value := range devices {
		if value.Name == dto.Name {
			return models.Device{}, errors.New("Device name must be unique per user")
		}
	}
	var device models.Device
	if dto.Type == 1 {
		device = res.airConditionerService.Add(dto).ToDevice()
	} else if dto.Type == 8 {
		device = res.evChargerService.Add(dto).ToDevice()
	} else if dto.Type == 7 {
		device = res.homeBatteryService.Add(dto).ToDevice()
		// todo add new case after adding new Device Class
	} else {
		device = dto.ToDevice()
		query := "INSERT INTO device (Name, Type, Picture, RealEstate, IsOnline)" +
			"VALUES ( ?, ?, ?, ?, ?);"
		result, err := res.db.Exec(query, device.Name, device.Type, device.Picture, device.RealEstate,
			device.IsOnline)
		if CheckIfError(err) {
			return models.Device{}, err
		}
		id, err := result.LastInsertId()
		device.Id = int(id)
	}

	res.mqtt.Publish(mqtt_client.TopicNewDevice+strconv.Itoa(device.Id), "new device created")
	return device, nil
}

func (ds *DeviceServiceImpl) getDevicesByUserID(userID int) ([]models.Device, error) {
	// Perform a database query to get devices by user ID
	rows, err := ds.db.Query("SELECT id, name, type, picture, realestate, isonline FROM device WHERE realestate IN (SELECT id FROM realestate WHERE userid = ?)", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the query result and populate the devices slice
	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.Picture, &device.RealEstate, &device.IsOnline)
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

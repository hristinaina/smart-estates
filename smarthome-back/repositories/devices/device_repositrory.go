package repositories

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"smarthome-back/dtos"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices"
	"smarthome-back/repositories"
	"strconv"
	"time"
)

type DeviceRepository interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	GetAll() []models.Device
	GetDevicesByUserID(userID int) ([]models.Device, error)
	Update(device models.Device) bool
	UpdateLastValue(id int, value float32) (bool, error)
	GetConsumptionDevicesByEstateId(userID int) ([]models.ConsumptionDevice, error)
	GetConsumptionDeviceDto(id int) (dtos.ConsumptionDeviceDto, error)
	GetConsumptionDevice(id int) (models.ConsumptionDevice, error)
	GetAvailability(dto dtos.ActionGraphRequest, isOnline string) []time.Time
}

type DeviceRepositoryImpl struct {
	db       *sql.DB
	influxDb influxdb2.Client
}

func NewDeviceRepository(db *sql.DB, influx influxdb2.Client) DeviceRepository {
	return &DeviceRepositoryImpl{db: db, influxDb: influx}
}

func (res *DeviceRepositoryImpl) GetAll() []models.Device {
	query := "SELECT * FROM device"
	rows, err := res.db.Query(query)
	if repositories.CheckIfError(err) {
		return nil
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate,
			&device.IsOnline, &device.StatusTimeStamp, &device.LastValue); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Device{}
		}
		devices = append(devices, device)
	}

	return devices
}

func (res *DeviceRepositoryImpl) GetConsumptionDevice(id int) (models.ConsumptionDevice, error) {
	query := `
		SELECT
			d.id,
			d.name,
			d.realEstate,
			d.isOnline,
			cd.powerSupply,
			cd.powerConsumption
		FROM
			device d
		JOIN
			consumptionDevice cd ON d.id = cd.deviceId
		WHERE
			d.id = ?
	`

	rows, err := res.db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	row := res.db.QueryRow(query, id)

	var cd models.ConsumptionDevice
	var device models.Device

	err = row.Scan(
		&device.Id,
		&device.Name,
		&device.RealEstate,
		&device.IsOnline,
		&cd.PowerSupply,
		&cd.PowerConsumption,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No consumption device found with the specified ID")
		} else {
			fmt.Println("Error retrieving solar panel:", err)
		}
		return models.ConsumptionDevice{}, err
	}
	cd.Device = device
	return cd, nil
}

func (res *DeviceRepositoryImpl) GetConsumptionDevicesByEstateId(id int) ([]models.ConsumptionDevice, error) {
	query := `
		SELECT
			d.id,
			d.name,
			d.realEstate,
			d.isOnline,
			cd.powerSupply,
			cd.powerConsumption
		FROM
			device d
		JOIN
			consumptionDevice cd ON d.id = cd.deviceId
		WHERE
			d.realEstate = ?
	`

	rows, err := res.db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set
	var consumptionDevices []models.ConsumptionDevice
	for rows.Next() {
		var device models.Device
		var cd models.ConsumptionDevice

		//todo da li treba da scan bude skroz ispunjen?
		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.RealEstate,
			&device.IsOnline,
			&cd.PowerSupply,
			&cd.PowerConsumption,
		)
		if err != nil {
			log.Fatal(err)
		}

		cd.Device = device
		consumptionDevices = append(consumptionDevices, cd)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return consumptionDevices, nil
}

func (res *DeviceRepositoryImpl) GetAllByEstateId(estateId int) []models.Device {
	query := "SELECT * FROM device WHERE REALESTATE = ?"
	rows, err := res.db.Query(query, estateId)
	if repositories.CheckIfError(err) {
		//todo raise an exception and catch it in controller?
		return nil
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		if err := rows.Scan(&device.Id, &device.Name, &device.Type,
			&device.RealEstate, &device.IsOnline, &device.StatusTimeStamp, &device.LastValue); err != nil {
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

	if repositories.CheckIfError(err) {
		return models.Device{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var (
			device models.Device
		)
		if err := rows.Scan(&device.Id, &device.Name, &device.Type,
			&device.RealEstate, &device.IsOnline, &device.StatusTimeStamp, &device.LastValue); err != nil {
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
	_, err := res.db.Exec(query, device.Name, device.Type, device.RealEstate, device.IsOnline, device.StatusTimeStamp,
		device.Id)
	if err != nil {
		fmt.Println("Failed to update device:", err)
		return false
	}
	return true
}

func (res *DeviceRepositoryImpl) UpdateLastValue(id int, value float32) (bool, error) {
	query := `UPDATE device
              SET device.LastValue = ? 
              WHERE Device.Id = ?`
	_, err := res.db.Exec(query, value, id)
	if repositories.CheckIfError(err) {
		return false, err
	}
	return true, nil
}

func (res *DeviceRepositoryImpl) GetConsumptionDeviceDto(id int) (dtos.ConsumptionDeviceDto, error) {
	query := `SELECT  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption
			  FROM ConsumptionDevice 
   			  WHERE ConsumptionDevice.DeviceId = ?`
	rows, err := res.db.Query(query, id)
	if repositories.IsError(err) {
		return dtos.ConsumptionDeviceDto{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Database connection closing error: ", err)
		}
	}(rows)

	var consumptionDevices []dtos.ConsumptionDeviceDto
	for rows.Next() {
		var (
			consumptionDevice dtos.ConsumptionDeviceDto
		)
		if err := rows.Scan(&consumptionDevice.PowerSupply, &consumptionDevice.PowerConsumption); err != nil {
			fmt.Println("Error: ", err.Error())
		}
		consumptionDevices = append(consumptionDevices, consumptionDevice)
	}

	if len(consumptionDevices) > 0 {
		return consumptionDevices[0], nil
	}
	// TODO: check if here should be Autonomous power supply
	return dtos.ConsumptionDeviceDto{PowerSupply: enumerations.Autonomous, PowerConsumption: 0}, nil
}

func (res *DeviceRepositoryImpl) GetAvailability(dto dtos.ActionGraphRequest, isOnline string) []time.Time {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	var end string
	if dto.EndDate == "-1" {
		end = "now()"
	} else {
		end = dto.EndDate
	}

	// Create InfluxDB query API
	queryAPI := res.influxDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(
		`from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r["_measurement"] == "device_status" and r["device_id"] == "%s" and r["_field"] == "status" and r["_value"] == %s)
		  |> map(fn: (r) => ({
			time: r["_time"]
		  }))`,
		influxBucket, dto.StartDate, end, strconv.Itoa(dto.DeviceId), isOnline,
	)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return nil
	}

	defer result.Close()
	var times []time.Time
	for result.Next() {
		times = append(times, result.Record().ValueByKey("time").(time.Time))
		if err := result.Err(); err != nil {
			fmt.Println("ERROR happened")
		}
	}
	return times
}

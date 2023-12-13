package repositories

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	models "smarthome-back/models/devices"
	devices "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
)

type LampRepository interface {
	Get(id int) (devices.Lamp, error)
	GetAll() ([]devices.Lamp, error)
	UpdateIsOnState(id int, isOn bool) (bool, error)
	UpdateLightningState(id int, lightningState int) (bool, error)
	GetLampData(from, to string) *api.QueryTableResult
}

type LampRepositoryImpl struct {
	db       *sql.DB
	influxdb influxdb2.Client
}

func NewLampRepository(db *sql.DB, influxdb influxdb2.Client) LampRepository {
	return &LampRepositoryImpl{db: db, influxdb: influxdb}
}

func (rl *LampRepositoryImpl) Get(id int) (devices.Lamp, error) {
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, Lamp.IsOn, Lamp.LightningLevel
			  FROM Lamp 
    		  JOIN ConsumptionDevice ON Lamp.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
   			  WHERE Device.Id = ?`
	rows, err := rl.db.Query(query, id)
	if repositories.IsError(err) {
		return devices.Lamp{}, err
	}
	defer rows.Close()
	lamps, err := ScanRows(rows)
	lamp := lamps[0]
	return lamp, err
}

func (rl *LampRepositoryImpl) GetAll() ([]devices.Lamp, error) {
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, Lamp.IsOn, Lamp.LightningLevel
			  FROM Lamp 
    		  JOIN ConsumptionDevice ON Lamp.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id`
	rows, err := rl.db.Query(query)
	if repositories.IsError(err) {
		return nil, err
	}
	defer rows.Close()
	return ScanRows(rows)
}

func (rl *LampRepositoryImpl) UpdateIsOnState(id int, isOn bool) (bool, error) {
	query := `UPDATE Lamp
              JOIN ConsumptionDevice ON Lamp.DeviceId = ConsumptionDevice.DeviceId
        	  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
              SET Lamp.IsOn = ? 
              WHERE Device.Id = ?`
	_, err := rl.db.Exec(query, isOn, id)
	if repositories.IsError(err) {
		return false, err
	}
	return true, nil
}

func (rl *LampRepositoryImpl) UpdateLightningState(id int, lightningState int) (bool, error) {
	query := `UPDATE Lamp
              JOIN ConsumptionDevice ON Lamp.DeviceId = ConsumptionDevice.DeviceId
        	  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
              SET Lamp.LightningLevel = ? 
              WHERE Device.Id = ?`
	_, err := rl.db.Exec(query, lightningState, id)
	if repositories.CheckIfError(err) {
		return false, err
	}
	return true, nil
}

func (rl *LampRepositoryImpl) GetLampData(from, to string) *api.QueryTableResult {
	client := rl.influxdb
	queryAPI := client.QueryAPI("Smart Home")
	// we are printing data that came in the last 10 minutes
	query := fmt.Sprintf(`from(bucket: "bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "measurement1")`, from, to)
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

// ScanRows mapping returned value from db to model - in this case in lamp model
func ScanRows(rows *sql.Rows) ([]devices.Lamp, error) {
	var lamps []devices.Lamp
	for rows.Next() {
		var (
			device     models.Device
			consDevice models.ConsumptionDevice
			lamp       devices.Lamp
		)

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate,
			&device.IsOnline, &consDevice.PowerSupply, &consDevice.PowerConsumption, &lamp.IsOn, &lamp.LightningLevel); err != nil {
			fmt.Println("Error: ", err.Error())
			return []devices.Lamp{}, err
		}
		consDevice.Device = device
		lamp.ConsumptionDevice = consDevice
		lamps = append(lamps, lamp)
	}

	return lamps, nil
}

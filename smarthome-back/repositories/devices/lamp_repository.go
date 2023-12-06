package repositories

import (
	"database/sql"
	"fmt"
	models "smarthome-back/models/devices"
	devices "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
)

type LampRepository struct {
	db *sql.DB
}

func NewLampRepository(db *sql.DB) *LampRepository {
	return &LampRepository{db: db}
}

func (rl *LampRepository) Get(id int) (devices.Lamp, error) {
	query := "SELECT * FROM lamp WHERE ID = ?"
	rows, err := rl.db.Query(query, id)
	if repositories.IsError(err) {
		return devices.Lamp{}, err
	}
	defer rows.Close()
	lamps, err := ScanRows(rows)
	lamp := lamps[0]
	return lamp, err
}

func (rl *LampRepository) GetAll() ([]devices.Lamp, error) {
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.Picture, Device.RealEstate, Device.IsOnline,
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

func (rl *LampRepository) UpdateIsOnState(id int, isOn bool) (bool, error) {
	query := "UPDATE lamp SET IsOn = ? WHERE Id = ?"
	_, err := rl.db.Exec(query, isOn, id)
	if repositories.IsError(err) {
		return false, err
	}
	return true, nil
}

func (rl *LampRepository) UpdateLightningState(id int, lightningState int) (bool, error) {
	query := "SELECT * FROM lamp SET LightningLevel = ? WHERE Id = ?"
	_, err := rl.db.Exec(query, lightningState, id)
	if repositories.CheckIfError(err) {
		return false, err
	}
	return true, nil
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

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.Picture, &device.RealEstate,
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

package repositories

import (
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
)

type SprinklerRepository interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() ([]models.Sprinkler, error)
	UpdateIsOn(isOn bool) (bool, error)
	Delete(id int) (bool, error)
	AddSpecialMode(mode models.SprinklerSpecialMode) (models.Sprinkler, error)
}

type SprinklerRepositoryImpl struct {
	db     *sql.DB
	influx influxdb2.Client
}

func NewSprinklerRepository(db *sql.DB, influx influxdb2.Client) SprinklerRepository {
	return &SprinklerRepositoryImpl{db: db, influx: influx}
}

func (repo *SprinklerRepositoryImpl) Get(id int) (models.Sprinkler, error) {
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, s.IsOn
			  FROM Sprinkler s 
    		  JOIN ConsumptionDevice ON s.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
   			  WHERE Device.Id = ? `

	rows, err := repo.db.Query(query, id)
	if repositories.IsError(err) {
		return models.Sprinkler{}, err
	}
	defer rows.Close()

	sprinklers, err := repo.scanRows(rows)
	if repositories.IsError(err) {
		return models.Sprinkler{}, err
	}
	sprinkler := sprinklers[0]
	// TODO: add here modes
	return sprinkler, nil

}

func (repo *SprinklerRepositoryImpl) GetAll() ([]models.Sprinkler, error) {
	return nil, nil
}

func (repo *SprinklerRepositoryImpl) UpdateIsOn(isOn bool) (bool, error) {
	return false, nil
}

func (repo *SprinklerRepositoryImpl) Delete(id int) (bool, error) {
	return false, nil
}

func (repo *SprinklerRepositoryImpl) AddSpecialMode(mode models.SprinklerSpecialMode) (models.Sprinkler, error) {
	return models.Sprinkler{}, nil
}

func (repo *SprinklerRepositoryImpl) scanRows(rows *sql.Rows) ([]models.Sprinkler, error) {
	var sprinklers []models.Sprinkler
	for rows.Next() {
		var (
			device     models2.Device
			consDevice models2.ConsumptionDevice
			sprinkler  models.Sprinkler
		)

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate,
			&device.IsOnline, &consDevice.PowerSupply, &consDevice.PowerConsumption, &sprinkler.IsOn); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.Sprinkler{}, err
		}
		consDevice.Device = device
		sprinkler.ConsumptionDevice = consDevice
		sprinklers = append(sprinklers, sprinkler)
	}

	return sprinklers, nil
}

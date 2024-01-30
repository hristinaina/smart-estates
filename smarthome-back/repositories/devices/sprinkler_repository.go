package repositories

import (
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
)

type SprinklerRepository interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() ([]models.Sprinkler, error)
	UpdateIsOn(id int, isOn bool) (bool, error)
	Delete(id int) (bool, error)
	AddSpecialMode(id int, mode models.SprinklerSpecialMode) (models.Sprinkler, error)
	Add(sprinkler models.Sprinkler) (models.Sprinkler, error)
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
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, s.IsOn
			  FROM Sprinkler s 
    		  JOIN ConsumptionDevice ON s.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id`

	rows, err := repo.db.Query(query)
	if repositories.IsError(err) {
		return nil, err
	}
	defer rows.Close()

	sprinklers, err := repo.scanRows(rows)
	if err != nil {
		return nil, err
	}
	return sprinklers, nil
}

func (repo *SprinklerRepositoryImpl) UpdateIsOn(id int, isOn bool) (bool, error) {
	query := `UPDATE Sprinkler s
			  JOIN ConsumptionDevice cd ON s.DeviceId = cd.DeviceId
			  JOIN Device d ON cd.DeviceId = d.Id
			  SET s.IsOn = ?
			  WHERE d.Id = ?`
	_, err := repo.db.Query(query, isOn, id)
	if repositories.IsError(err) {
		return false, err
	}
	return true, nil
}

func (repo *SprinklerRepositoryImpl) Delete(id int) (bool, error) {
	_, err := repo.Get(id)
	if err != nil {
		return false, err
	}

	tx, err := repo.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	_, err = tx.Exec("DELETE FROM SprinklerSpecialMode WHERE DeviceId = ?", id)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("DELETE FROM Sprinkler WHERE DeviceId = ?", id)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("DELETE FROM ConsumptionDevice WHERE DeviceId = ?", id)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("DELETE FROM Device WHERE Id = ?", id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *SprinklerRepositoryImpl) AddSpecialMode(id int, mode models.SprinklerSpecialMode) (models.Sprinkler, error) {
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

func (repo *SprinklerRepositoryImpl) Add(device models.Sprinkler) (models.Sprinkler, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return models.Sprinkler{}, err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.ConsumptionDevice.Device.Name, device.ConsumptionDevice.Device.Type,
		device.ConsumptionDevice.Device.RealEstate, device.ConsumptionDevice.Device.IsOnline)
	if err != nil {
		return models.Sprinkler{}, err
	}

	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.Sprinkler{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO ConsumptionDevice(DeviceId, PowerSupply, PowerConsumption)
							VALUES (?, ?, ?)`, deviceID, device.ConsumptionDevice.PowerSupply,
		device.ConsumptionDevice.PowerConsumption)
	if err != nil {
		return models.Sprinkler{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO Sprinkler(DeviceId, IsOn)
							VALUES (?, ?)`, deviceID, device.IsOn)
	if err != nil {
		return models.Sprinkler{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.Sprinkler{}, err
	}
	device.ConsumptionDevice.Device.Id = int(deviceID)

	return device, nil
}

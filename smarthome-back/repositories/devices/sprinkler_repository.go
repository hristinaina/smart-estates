package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"smarthome-back/cache"
	"smarthome-back/enumerations"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SprinklerRepository interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() ([]models.Sprinkler, error)
	UpdateIsOn(id int, isOn bool) (bool, error)
	Delete(id int) (bool, error)
	AddSpecialMode(id int, mode models.SprinklerSpecialMode) (models.SprinklerSpecialMode, error)
	Add(sprinkler models.Sprinkler) (models.Sprinkler, error)
	GetSpecialModes(deviceId int) ([]models.SprinklerSpecialMode, error)
	DeleteSpecialMode(id int) (bool, error)
	GetSpecialMode(id int) (models.SprinklerSpecialMode, error)
}

type SprinklerRepositoryImpl struct {
	db           *sql.DB
	influx       influxdb2.Client
	cacheService *cache.CacheService
}

func NewSprinklerRepository(db *sql.DB, influx influxdb2.Client, cacheService cache.CacheService) SprinklerRepository {
	return &SprinklerRepositoryImpl{db: db, influx: influx, cacheService: &cacheService}
}

func (repo *SprinklerRepositoryImpl) SelectQuery(id int) (models.Sprinkler, error) {
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

	return sprinkler, err
}

func (repo *SprinklerRepositoryImpl) Get(id int) (models.Sprinkler, error) {
	cacheKey := fmt.Sprintf("sprinkler_%d", id)

	var sprinkler models.Sprinkler
	if found, err := repo.cacheService.GetFromCache(cacheKey, &sprinkler); found {
		return sprinkler, err
	}

	sprinkler, _ = repo.SelectQuery(id)
	// TODO: add here modes

	if err := repo.cacheService.SetToCache(cacheKey, sprinkler); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

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

	sprinkler, _ := repo.SelectQuery(id)

	cacheKey := fmt.Sprintf("sprinkler_%d", id)
	if err := repo.cacheService.SetToCache(cacheKey, sprinkler); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
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

func (repo *SprinklerRepositoryImpl) AddSpecialMode(id int, mode models.SprinklerSpecialMode) (models.SprinklerSpecialMode, error) {
	selectedDays := mode.SelectedDaysToString()
	query := `INSERT INTO SprinklerSpecialMode (DeviceId, StartTime, EndTime, SelectedDays) VALUES (?, ?, ?, ?)`
	res, err := repo.db.Exec(query, id, mode.StartTime, mode.EndTime, selectedDays)
	if repositories.CheckIfError(err) {
		return models.SprinklerSpecialMode{}, err
	}
	modeId, err := res.LastInsertId()
	mode.DeviceId = id
	mode.Id = int(modeId)
	return mode, nil
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

	cacheKey := fmt.Sprintf("sprinkler_%d", device.ConsumptionDevice.Device.Id)
	if err := repo.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = repo.cacheService.AddDevicesByRealEstate(device.ConsumptionDevice.Device.RealEstate, device.ConsumptionDevice.Device)
	return device, nil
}

func (repo *SprinklerRepositoryImpl) GetSpecialModes(id int) ([]models.SprinklerSpecialMode, error) {
	query := `SELECT Id, DeviceId, StartTime, EndTime, SelectedDays
			  FROM SprinklerSpecialMode
              WHERE DeviceId = ?`
	rows, err := repo.db.Query(query, id)
	if repositories.IsError(err) {
		return nil, err
	}
	defer rows.Close()

	modes, err := repo.scanModeRows(rows)
	if repositories.IsError(err) {
		return nil, err
	}
	return modes, nil

}

func (repo *SprinklerRepositoryImpl) DeleteSpecialMode(id int) (bool, error) {
	_, err := repo.GetSpecialMode(id)
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

	_, err = tx.Exec("DELETE FROM SprinklerSpecialMode WHERE Id = ?", id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *SprinklerRepositoryImpl) GetSpecialMode(id int) (models.SprinklerSpecialMode, error) {
	query := `SELECT * 
			   FROM SprinklerSpecialMode
			   WHERE Id = ?`

	rows, err := repo.db.Query(query, id)
	if err != nil {
		return models.SprinklerSpecialMode{}, err
	}
	defer rows.Close()

	modes, err := repo.scanModeRows(rows)
	if repositories.IsError(err) {
		return models.SprinklerSpecialMode{}, err
	}
	mode := modes[0]
	return mode, nil
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

func (repo *SprinklerRepositoryImpl) scanModeRows(rows *sql.Rows) ([]models.SprinklerSpecialMode, error) {
	var modes []models.SprinklerSpecialMode
	for rows.Next() {
		var (
			mode         models.SprinklerSpecialMode
			selectedDays string
		)

		if err := rows.Scan(&mode.Id, &mode.DeviceId, &mode.StartTime, &mode.EndTime, &selectedDays); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.SprinklerSpecialMode{}, err
		}
		days, err := enumerations.ConvertStringsToEnumValues(selectedDays)
		if err != nil {
			return []models.SprinklerSpecialMode{}, err
		}
		mode.SelectedDays = days
		modes = append(modes, mode)
	}
	return modes, nil
}

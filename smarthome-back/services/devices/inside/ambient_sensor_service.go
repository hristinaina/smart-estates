package inside

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
)

type AmbientSensorService interface {
	Get(id int) (models.ConsumptionDevice, error)
	Add(dto dtos.DeviceDTO) models.ConsumptionDevice
}

type AmbientSensorServiceImpl struct {
	db           *sql.DB
	cacheService cache.CacheService
}

func NewAmbientSensorService(db *sql.DB, cacheService *cache.CacheService) AmbientSensorService {
	return &AmbientSensorServiceImpl{db: db, cacheService: *cacheService}
}

func (as *AmbientSensorServiceImpl) Get(id int) (models.ConsumptionDevice, error) {
	cacheKey := fmt.Sprintf("as_%d", id)

	var sensor models.ConsumptionDevice
	if found, err := as.cacheService.GetFromCache(cacheKey, &sensor); found {
		return sensor, err
	}

	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			ConsumptionDevice.PowerSupply,
			ConsumptionDevice.PowerConsumption
		FROM
			ConsumptionDevice
		JOIN 
			Device ON ConsumptionDevice.DeviceId = Device.Id
		WHERE
			Device.Id = ?;
	`
	// Execute the query
	rows, err := as.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return models.ConsumptionDevice{}, err
	}
	defer rows.Close()

	var device models.Device
	var consDevice models.ConsumptionDevice

	for rows.Next() {
		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Type,
			&device.RealEstate,
			&device.IsOnline,
			&device.StatusTimeStamp,
			&consDevice.PowerSupply,
			&consDevice.PowerConsumption,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return models.ConsumptionDevice{}, err
		}
		consDevice.Device = device
	}

	if err := as.cacheService.SetToCache(cacheKey, consDevice); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return consDevice, nil
}

func (as *AmbientSensorServiceImpl) Add(dto dtos.DeviceDTO) models.ConsumptionDevice {
	device := dto.ToAmbientSensor()
	tx, err := as.db.Begin()
	if err != nil {
		return models.ConsumptionDevice{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.PowerSupply, device.PowerConsumption)
	if err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return models.ConsumptionDevice{}
	}

	device.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("as_%d", device.Device.Id)
	if err := as.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	err = as.cacheService.AddDevicesByRealEstate(device.Device.RealEstate, device.Device)

	return device
}

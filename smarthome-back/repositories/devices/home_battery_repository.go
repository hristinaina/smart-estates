package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/energetic"
)

type HomeBatteryRepository interface {
	Add(estate dtos.DeviceDTO) energetic.HomeBattery
	GetAllByEstateId(id int) ([]energetic.HomeBattery, error)
	Update(device energetic.HomeBattery) bool
	Get(id int) energetic.HomeBattery
}

type HomeBatteryRepositoryImpl struct {
	db           *sql.DB
	cacheService cache.CacheService
}

func NewHomeBatteryRepository(db *sql.DB, cacheService cache.CacheService) HomeBatteryRepository {
	return &HomeBatteryRepositoryImpl{db: db, cacheService: cacheService}
}

func (s *HomeBatteryRepositoryImpl) GetAllByEstateId(id int) ([]energetic.HomeBattery, error) {
	cacheKey := fmt.Sprintf("battery_estate_%d", id)

	var batteries []energetic.HomeBattery
	if found, err := s.cacheService.GetFromCache(cacheKey, &batteries); found {
		return batteries, err
	}

	query := `
		SELECT
			d.id,
			d.name,
			d.realEstate,
			d.isOnline,
			hb.size,
			hb.currentValue
		FROM
			device d
		JOIN
			homeBattery hb ON d.id = hb.deviceId
		WHERE
			d.realEstate = ?
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var device models.Device
		var hb energetic.HomeBattery

		//todo da li treba da scan bude skroz ispunjen?
		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.RealEstate,
			&device.IsOnline,
			&hb.Size,
			&hb.CurrentValue,
		)
		if err != nil {
			log.Fatal(err)
		}

		hb.Device = device
		batteries = append(batteries, hb)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if err := s.cacheService.SetToCache(cacheKey, batteries); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	return batteries, nil
}

func (s *HomeBatteryRepositoryImpl) Add(dto dtos.DeviceDTO) energetic.HomeBattery {
	// TODO: add some validation and exception throwing
	device := dto.ToHomeBattery()
	tx, err := s.db.Begin()
	if err != nil {
		return energetic.HomeBattery{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return energetic.HomeBattery{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return energetic.HomeBattery{}
	}

	// Insert the new Home Battery into the Home Battery table
	result, err = tx.Exec(`
		INSERT INTO HomeBattery (DeviceId, Size)
		VALUES (?, ?)
	`, deviceID, device.Size)
	if err != nil {
		return energetic.HomeBattery{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return energetic.HomeBattery{}
	}
	device.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("battery_estate_%d", device.Device.Id)
	if err := s.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = s.cacheService.AddDevicesByRealEstate(device.Device.RealEstate, device.Device)
	return device
}

func (res *HomeBatteryRepositoryImpl) Update(device energetic.HomeBattery) bool {
	query := "UPDATE homeBattery SET currentValue = ? WHERE deviceId = ?"
	_, err := res.db.Exec(query, device.CurrentValue, device.Device.Id)
	if err != nil {
		fmt.Println("Failed to update device:", err)
		return false
	}

	cacheKey := fmt.Sprintf("battery_estate_%d", device.Device.Id)
	if err := res.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	return true
}

func (res *HomeBatteryRepositoryImpl) Get(id int) energetic.HomeBattery {
	cacheKey := fmt.Sprintf("battery_%d", id)

	var hb energetic.HomeBattery
	if found, _ := res.cacheService.GetFromCache(cacheKey, &hb); found {
		return hb
	}

	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			HomeBattery.Size,
			HomeBattery.CurrentValue
		FROM
			HomeBattery
		JOIN Device ON HomeBattery.DeviceId = Device.Id
		WHERE
			Device.Id = ?
	`

	// Execute the query
	row := res.db.QueryRow(query, id)

	var device models.Device

	err := row.Scan(
		&device.Id,
		&device.Name,
		&device.Type,
		&device.RealEstate,
		&device.IsOnline,
		&device.StatusTimeStamp,
		&hb.Size,
		&hb.CurrentValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No solar panel found with the specified ID")
		} else {
			fmt.Println("Error retrieving solar panel:", err)
		}
		return energetic.HomeBattery{}
	}
	hb.Device = device

	if err := res.cacheService.SetToCache(cacheKey, hb); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return hb
}

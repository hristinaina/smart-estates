package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/energetic"
)

type EVChargerRepository interface {
	Get(id int) energetic.EVCharger
	Add(dto dtos.DeviceDTO) energetic.EVCharger
}

type EVChargerRepositoryImpl struct {
	db           *sql.DB
	cacheService *cache.CacheService
}

func NewEVChargerRepository(db *sql.DB, cacheService cache.CacheService) EVChargerRepository {
	return &EVChargerRepositoryImpl{db: db, cacheService: &cacheService}
}

func (s *EVChargerRepositoryImpl) Get(id int) energetic.EVCharger {
	cacheKey := fmt.Sprintf("charger_%d", id)

	var sp energetic.EVCharger
	if found, _ := s.cacheService.GetFromCache(cacheKey, &sp); found {
		return sp
	}

	query := `SELECT
				d.id,
				d.name,
				d.type,
				d.realEstate,
				d.isOnline,
				hb.chargingPower,
				hb.connections
			FROM
				device d
			JOIN
				evCharger hb ON d.id = hb.deviceId
			WHERE
				d.id = ?`

	// Execute the query
	row := s.db.QueryRow(query, id)

	var device models.Device

	err := row.Scan(
		&device.Id,
		&device.Name,
		&device.Type,
		&device.RealEstate,
		&device.IsOnline,
		&sp.ChargingPower,
		&sp.Connections,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No charger found with the specified ID")
		} else {
			fmt.Println("Error retrieving charger:", err)
		}
		return energetic.EVCharger{}
	}
	sp.Device = device

	if err := s.cacheService.SetToCache(cacheKey, sp); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return sp
}

func (s *EVChargerRepositoryImpl) Add(dto dtos.DeviceDTO) energetic.EVCharger {
	device := dto.ToEVCharger()
	tx, err := s.db.Begin()
	if err != nil {
		return energetic.EVCharger{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return energetic.EVCharger{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return energetic.EVCharger{}
	}

	// Insert the new EVCharger into the EVCharger table
	result, err = tx.Exec(`
		INSERT INTO EvCharger (DeviceId, ChargingPower, Connections)
		VALUES (?, ?, ?)
	`, deviceID, device.ChargingPower, device.Connections)
	if err != nil {
		return energetic.EVCharger{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return energetic.EVCharger{}
	}
	device.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("charger_%d", device.Device.Id)
	if err := s.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = s.cacheService.AddDevicesByRealEstate(device.Device.RealEstate, device.Device)
	return device
}

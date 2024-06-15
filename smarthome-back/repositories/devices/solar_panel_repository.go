package repositories

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/energetic"
)

type SolarPanelRepository interface {
	Get(id int) energetic.SolarPanel
	UpdateSP(device energetic.SolarPanel) bool
	Add(dto dtos.DeviceDTO) energetic.SolarPanel
}

type SolarPanelRepositoryImpl struct {
	db           *sql.DB
	cacheService *cache.CacheService
}

func NewSolarPanelRepository(db *sql.DB, cacheService cache.CacheService) SolarPanelRepository {
	return &SolarPanelRepositoryImpl{db: db, cacheService: &cacheService}
}

func (s *SolarPanelRepositoryImpl) Get(id int) energetic.SolarPanel {
	cacheKey := fmt.Sprintf("sp_%d", id)

	var sp energetic.SolarPanel
	if found, _ := s.cacheService.GetFromCache(cacheKey, &sp); found {
		return sp
	}

	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			SolarPanel.SurfaceArea,
			SolarPanel.Efficiency,
			SolarPanel.NumberOfPanels,
			SolarPanel.IsOn
		FROM
			SolarPanel
		JOIN Device ON SolarPanel.DeviceId = Device.Id
		WHERE
			Device.Id = ?
	`

	// Execute the query
	row := s.db.QueryRow(query, id)

	var device models.Device

	err := row.Scan(
		&device.Id,
		&device.Name,
		&device.Type,
		&device.RealEstate,
		&device.IsOnline,
		&device.StatusTimeStamp,
		&sp.SurfaceArea,
		&sp.Efficiency,
		&sp.NumberOfPanels,
		&sp.IsOn,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No solar panel found with the specified ID")
		} else {
			fmt.Println("Error retrieving solar panel:", err)
		}
		return energetic.SolarPanel{}
	}
	sp.Device = device

	if err := s.cacheService.SetToCache(cacheKey, sp); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return sp
}

func (s *SolarPanelRepositoryImpl) Add(dto dtos.DeviceDTO) energetic.SolarPanel {
	// TODO: add some validation and exception throwing
	device := dto.ToSolarPanel()
	tx, err := s.db.Begin()
	if err != nil {
		return energetic.SolarPanel{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return energetic.SolarPanel{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return energetic.SolarPanel{}
	}

	// Insert the new SolarPanel into the SolarPanel table
	result, err = tx.Exec(`
		INSERT INTO SolarPanel (DeviceId, SurfaceArea, Efficiency, NumberOfPanels, IsOn)
		VALUES (?, ?, ?, ?, ?)
	`, deviceID, device.SurfaceArea, device.Efficiency, device.NumberOfPanels, false)
	if err != nil {
		return energetic.SolarPanel{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return energetic.SolarPanel{}
	}
	device.Device.Id = int(deviceID)
	return device
}

func (s *SolarPanelRepositoryImpl) UpdateSP(device energetic.SolarPanel) bool {
	query := "UPDATE solarPanel SET surfaceArea = ?, efficiency = ?, numberOfPanels = ?, isOn = ? WHERE deviceId = ?"
	_, err := s.db.Exec(query, device.SurfaceArea, device.Efficiency, device.NumberOfPanels, device.IsOn, device.Device.Id)
	if err != nil {
		fmt.Println("Failed to update device:", err)
		return false
	}
	return true
}

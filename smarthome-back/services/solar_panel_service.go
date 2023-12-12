package services

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
)

type SolarPanelService interface {
	Add(estate dto.DeviceDTO) models.SolarPanel
	Get(id int) models.SolarPanel
}

type SolarPanelServiceImpl struct {
	db *sql.DB
}

func NewSolarPanelService(db *sql.DB) SolarPanelService {
	return &SolarPanelServiceImpl{db: db}
}

func (s *SolarPanelServiceImpl) Get(id int) models.SolarPanel {
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
			SolarPanel.IsOn
		FROM
			SolarPanel
		JOIN Device ON SolarPanel.DeviceId = Device.Id
		WHERE
			Device.Id = ?
	`

	// Execute the query
	row := s.db.QueryRow(query, id)

	var sp models.SolarPanel
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
		&sp.IsOn,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No solar panel found with the specified ID")
		} else {
			fmt.Println("Error retrieving solar panel:", err)
		}
		return models.SolarPanel{}
	}
	sp.Device = device
	return sp
}

func (s *SolarPanelServiceImpl) Add(dto dto.DeviceDTO) models.SolarPanel {
	// TODO: add some validation and exception throwing
	device := dto.ToSolarPanel()
	tx, err := s.db.Begin()
	if err != nil {
		return models.SolarPanel{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Name, device.Device.Type, device.Device.RealEstate,
		device.Device.IsOnline)
	if err != nil {
		return models.SolarPanel{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.SolarPanel{}
	}

	// Insert the new SolarPanel into the SolarPanel table
	result, err = tx.Exec(`
		INSERT INTO SolarPanel (DeviceId, SurfaceArea, Efficiency, IsOn)
		VALUES (?, ?, ?, ?)
	`, deviceID, device.SurfaceArea, device.Efficiency, false)
	if err != nil {
		return models.SolarPanel{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.SolarPanel{}
	}
	device.Device.Id = int(deviceID)
	return device
}

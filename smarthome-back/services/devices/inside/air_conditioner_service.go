package inside

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/inside"
	"strings"
)

type AirConditionerService interface {
	Add(estate dtos.DeviceDTO) inside.AirConditioner
	Get(id int) inside.AirConditioner
	UpdateSpecialMode(deviceId int, mode string, startTime string, endTime string, temperature float32, selectedDays string) error
	AddSpecialModes(deviceID int, mode string, startTime string, endTime string, temperature float32, selectedDays string) error
	DeleteSpecialMode(deviceID int, smDTO []dtos.SpecialModeDTO) error
}

type AirConditionerServiceImpl struct {
	db           *sql.DB
	cacheService *cache.CacheService
}

func NewAirConditionerService(db *sql.DB, cacheService *cache.CacheService) AirConditionerService {
	return &AirConditionerServiceImpl{db: db, cacheService: cacheService}
}

func (s *AirConditionerServiceImpl) selectQuery(id int) inside.AirConditioner {
	var ac inside.AirConditioner

	query := `
		SELECT
			Device.Id,
			Device.Name,
			Device.Type,
			Device.RealEstate,
			Device.IsOnline,
			Device.StatusTimeStamp,
			ConsumptionDevice.PowerSupply,
			ConsumptionDevice.PowerConsumption,
			AirConditioner.MinTemperature,
			AirConditioner.MaxTemperature,
			AirConditioner.Mode,
			SpecialModes.StartTime,
			SpecialModes.EndTime,
			SpecialModes.Mode,
			SpecialModes.Temperature,
			SpecialModes.SelectedDays
		FROM
			AirConditioner
		JOIN ConsumptionDevice ON AirConditioner.DeviceId = ConsumptionDevice.DeviceId
		JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
		LEFT JOIN SpecialModes ON AirConditioner.DeviceId = SpecialModes.DeviceId
		WHERE
			Device.Id = ?
	`

	// Execute the query
	rows, err := s.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return inside.AirConditioner{}
	}
	defer rows.Close()

	var device models.Device
	var consDevice models.ConsumptionDevice
	var specialModes []inside.SpecialMode

	for rows.Next() {
		var startTimeStr, endTimeStr sql.NullString
		var mode sql.NullString
		var selectedDays sql.NullString
		var temperature sql.NullFloat64

		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Type,
			&device.RealEstate,
			&device.IsOnline,
			&device.StatusTimeStamp,
			&consDevice.PowerSupply,
			&consDevice.PowerConsumption,
			&ac.MinTemperature,
			&ac.MaxTemperature,
			&ac.Mode,
			&startTimeStr,
			&endTimeStr,
			&mode,
			&temperature,
			&selectedDays,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return inside.AirConditioner{}
		}

		consDevice.Device = device
		ac.Device = consDevice

		specialMode := inside.SpecialMode{}
		if startTimeStr.Valid {
			specialMode.StartTime = startTimeStr.String
		}
		if endTimeStr.Valid {
			specialMode.EndTime = endTimeStr.String
		}
		if mode.Valid {
			specialMode.Mode = mode.String
		}
		if temperature.Valid {
			specialMode.Temperature = float32(temperature.Float64)
		}
		if selectedDays.Valid {
			specialMode.SelectedDays = selectedDays.String
		}
		specialModes = append(specialModes, specialMode)
	}

	ac.SpecialMode = specialModes

	return ac
}

func (s *AirConditionerServiceImpl) Get(id int) inside.AirConditioner {
	cacheKey := fmt.Sprintf("ac_%d", id)

	var ac inside.AirConditioner
	if found, _ := s.cacheService.GetFromCache(cacheKey, &ac); found {
		return ac
	}

	ac = s.selectQuery(id)

	if err := s.cacheService.SetToCache(cacheKey, ac); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return ac
}

func (s *AirConditionerServiceImpl) Add(dto dtos.DeviceDTO) inside.AirConditioner {
	device := dto.ToAirConditioner()
	tx, err := s.db.Begin()
	if err != nil {
		return inside.AirConditioner{}
	}
	defer tx.Rollback()

	// Insert the new device into the Device table
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.Device.Device.Name, device.Device.Device.Type, device.Device.Device.RealEstate,
		device.Device.Device.IsOnline)
	if err != nil {
		fmt.Println(err)
		return inside.AirConditioner{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return inside.AirConditioner{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.Device.PowerSupply, device.Device.PowerConsumption)
	if err != nil {
		fmt.Println(err)
		return inside.AirConditioner{}
	}

	// Insert the new air conditioner into the AirConditioner table
	result, err = tx.Exec(`
		INSERT INTO AirConditioner (DeviceId, MinTemperature, MaxTemperature, Mode)
		VALUES (?, ?, ?, ?)
	`, deviceID, device.MinTemperature, device.MaxTemperature, device.Mode)
	if err != nil {
		fmt.Println("ovde je greska")
		fmt.Println(err)
		return inside.AirConditioner{}
	}
	fmt.Println("OVDE SAAAAAAAAM")
	fmt.Println(device.SpecialMode)
	if len(device.SpecialMode) != 0 {
		fmt.Println(device.SpecialMode)
		for _, mode := range device.SpecialMode {
			result, err = tx.Exec(`
			INSERT INTO specialModes (DeviceId, StartTime, EndTime, Mode, Temperature, SelectedDays)
			VALUES (?, ?, ?, ?, ?, ?)
		`, deviceID, mode.StartTime, mode.EndTime, mode.Mode, mode.Temperature, mode.SelectedDays)
			if err != nil {
				fmt.Println("ovde jeeeeeeeeeee")
				fmt.Println(err)
				return inside.AirConditioner{}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return inside.AirConditioner{}
	}

	device.Device.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("ac_%d", device.Device.Device.Id)
	if err := s.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = s.cacheService.AddDevicesByRealEstate(device.Device.Device.RealEstate, device.Device.Device)

	return device
}

func (s *AirConditionerServiceImpl) UpdateSpecialMode(deviceId int, mode string, startTime string, endTime string, temperature float32, selectedDays string) error {
	query := `
        UPDATE specialModes
        SET Mode = ?, StartTime = ?, EndTime = ?, Temperature = ?, SelectedDays = ?
        WHERE DeviceId = ?
    `

	_, err := s.db.Exec(query, mode, startTime, endTime, temperature, selectedDays, deviceId)
	if err != nil {
		return fmt.Errorf("failed to update special mode: %v", err)
	}

	ac := s.selectQuery(deviceId)

	cacheKey := fmt.Sprintf("ac_%d", deviceId)
	if err := s.cacheService.SetToCache(cacheKey, ac); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return nil
}

type SpecialModeDB struct {
	ID           int
	DeviceID     int
	StartTime    string
	EndTime      string
	Mode         string
	Temperature  float32
	SelectedDays string
}

func (s *AirConditionerServiceImpl) getSpecialModesByDeviceId(deviceID int) ([]SpecialModeDB, error) {
	rows, err := s.db.Query("SELECT * FROM specialModes WHERE DeviceId = ?", deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduledTerms []SpecialModeDB

	for rows.Next() {
		var term SpecialModeDB
		err := rows.Scan(&term.ID, &term.DeviceID, &term.StartTime, &term.EndTime, &term.Mode, &term.Temperature, &term.SelectedDays)
		if err != nil {
			return nil, err
		}
		scheduledTerms = append(scheduledTerms, term)
	}

	return scheduledTerms, nil
}

func (s *AirConditionerServiceImpl) DeleteSpecialMode(deviceID int, smDTO []dtos.SpecialModeDTO) error {
	data, err := s.getSpecialModesByDeviceId(deviceID)
	if err != nil {
		fmt.Println("greska kod get")
		fmt.Println(err)
	}

	forDelete := true

	// ako postoji u bazi ali ne i u listi
	for _, bazaTerm := range data {
		for _, mode := range smDTO {
			if bazaTerm.DeviceID == deviceID && mode.Start == bazaTerm.StartTime && mode.End == bazaTerm.EndTime && mode.SelectedMode == bazaTerm.Mode && mode.Temperature == bazaTerm.Temperature && strings.Join(mode.SelectedDays, ",") == bazaTerm.SelectedDays {
				forDelete = false
			}
		}
		if forDelete {
			fmt.Println("obrisano")
			_, err := s.db.Exec("DELETE FROM specialModes WHERE DeviceId = ? AND StartTime = ? AND EndTime = ? AND Mode = ? AND Temperature = ? AND SelectedDays = ?", deviceID, bazaTerm.StartTime, bazaTerm.EndTime, bazaTerm.Mode, bazaTerm.Temperature, bazaTerm.SelectedDays)
			if err != nil {
				return err
			}
			break
		}
	}

	ac := s.selectQuery(deviceID)

	cacheKey := fmt.Sprintf("ac_%d", deviceID)
	if err := s.cacheService.SetToCache(cacheKey, ac); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return nil
}

func (s *AirConditionerServiceImpl) AddSpecialModes(deviceID int, mode string, startTime string, endTime string, temperature float32, selectedDays string) error {
	data, err := s.getSpecialModesByDeviceId(deviceID)
	if err != nil {
		fmt.Println("greska kod get")
		fmt.Println(err)
	}

	found := false

	// Provera da li termin postoji u bazi
	for _, bazaTerm := range data {
		if bazaTerm.DeviceID == deviceID && startTime == bazaTerm.StartTime && endTime == bazaTerm.EndTime && mode == bazaTerm.Mode && temperature == bazaTerm.Temperature && selectedDays == bazaTerm.SelectedDays {
			found = true
			break
		}
	}

	// Ako termin nije pronađen u bazi, dodajemo ga
	if !found {
		_, err := s.db.Exec("INSERT INTO specialModes (DeviceId, StartTime, EndTime, Mode, Temperature, SelectedDays) VALUES (?, ?, ?, ?, ?, ?)", deviceID, startTime, endTime, mode, temperature, selectedDays)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	ac := s.selectQuery(deviceID)

	cacheKey := fmt.Sprintf("ac_%d", deviceID)
	if err := s.cacheService.SetToCache(cacheKey, ac); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return nil
}

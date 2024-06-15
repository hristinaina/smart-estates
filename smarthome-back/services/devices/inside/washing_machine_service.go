package inside

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/inside"
	"strings"
	"time"
)

type WashingMachineService interface {
	Add(estate dtos.DeviceDTO) inside.WashingMachine
	Get(id int) inside.WashingMachine
	AddScheduledMode(deviceId, modeId int, startTime string) error
	GetAllScheduledModesForDevice(deviceId int) []inside.ScheduledMode
}

type WashingMachineServiceImpl struct {
	db           *sql.DB
	cacheService cache.CacheService
}

func NewWashingMachineService(db *sql.DB, cacheService *cache.CacheService) WashingMachineService {
	return &WashingMachineServiceImpl{db: db, cacheService: *cacheService}
}

func (s *WashingMachineServiceImpl) Get(id int) inside.WashingMachine {
	cacheKey := fmt.Sprintf("wm_%d", id)

	var wm inside.WashingMachine
	if found, _ := s.cacheService.GetFromCache(cacheKey, &wm); found {
		return wm
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
		ConsumptionDevice.PowerConsumption,
		WashingMachine.Mode,
		MachineScheduledMode.StartTime,
		MachineMode.Name AS ModeName,
		MachineMode.Duration,
		MachineMode.Temp
	FROM
		WashingMachine
	JOIN
		ConsumptionDevice ON WashingMachine.DeviceId = ConsumptionDevice.DeviceId
	JOIN
		Device ON ConsumptionDevice.DeviceId = Device.Id
	LEFT JOIN
		MachineScheduledMode ON WashingMachine.DeviceId = MachineScheduledMode.DeviceId
	LEFT JOIN
		MachineMode ON MachineScheduledMode.ModeId = MachineMode.Id
	WHERE
		Device.Id = ?;
	`

	// Execute the query
	rows, err := s.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return inside.WashingMachine{}
	}
	defer rows.Close()

	var device models.Device
	var consDevice models.ConsumptionDevice

	for rows.Next() {
		var startTime sql.NullString
		var name sql.NullString
		var duration sql.NullFloat64
		var temperature sql.NullString
		var modeNames sql.NullString

		err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Type,
			&device.RealEstate,
			&device.IsOnline,
			&device.StatusTimeStamp,
			&consDevice.PowerSupply,
			&consDevice.PowerConsumption,
			&modeNames,
			&startTime,
			&name,
			&duration,
			&temperature,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return inside.WashingMachine{}
		}

		consDevice.Device = device
		wm.Device = consDevice

		wm.ModeName = modeNames.String

		modes := s.findModeBasedOnName(modeNames.String)
		wm.Mode = modes
	}

	if err := s.cacheService.SetToCache(cacheKey, wm); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return wm
}

func (s *WashingMachineServiceImpl) findModeBasedOnName(names string) []inside.Mode {
	var modes []inside.Mode

	parts := strings.Split(names, ",")

	for _, part := range parts {
		var id int
		var duration int
		var temp string
		var modeName string

		query := "SELECT Id, Name, Duration, Temp FROM machineMode WHERE Name = ?"

		row := s.db.QueryRow(query, part)

		err := row.Scan(&id, &modeName, &duration, &temp)
		if err != nil {
			fmt.Println(err)
			return modes
		}

		mode := inside.Mode{
			Id:          id,
			Name:        modeName,
			Duration:    duration,
			Temperature: temp,
		}

		modes = append(modes, mode)
	}
	return modes
}

func (s *WashingMachineServiceImpl) Add(dto dtos.DeviceDTO) inside.WashingMachine {
	device := dto.ToWashingMachine()
	tx, err := s.db.Begin()
	if err != nil {
		return inside.WashingMachine{}
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
		return inside.WashingMachine{}
	}

	// Get the last inserted device ID
	deviceID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return inside.WashingMachine{}
	}

	// Insert the new consumption device into the ConsumptionDevice table
	_, err = tx.Exec(`
		INSERT INTO ConsumptionDevice (DeviceId, PowerSupply, PowerConsumption)
		VALUES (?, ?, ?)
	`, deviceID, device.Device.PowerSupply, device.Device.PowerConsumption)
	if err != nil {
		fmt.Println(err)
		return inside.WashingMachine{}
	}

	// Insert the new washing machine into the WashingMachine table
	result, err = tx.Exec(`
		INSERT INTO WashingMachine (DeviceId, Mode)
		VALUES (?, ?)
	`, deviceID, device.ModeName)
	if err != nil {
		fmt.Println(err)
		return inside.WashingMachine{}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return inside.WashingMachine{}
	}

	device.Device.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("wm_%d", device.Device.Device.Id)
	if err := s.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = s.cacheService.AddDevicesByRealEstate(device.Device.Device.RealEstate, device.Device.Device)
	return device
}

func (s *WashingMachineServiceImpl) AddScheduledMode(deviceId, modeId int, startTime string) error {
	query := "INSERT INTO machineScheduledMode (Id, DeviceId, StartTime, ModeId) VALUES (?, ?, ?, ?);"
	_, err := s.db.Exec(query, s.generateId(), deviceId, startTime, modeId)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Failed to save mode: %v", err)

	}
	return nil
}

func (res *WashingMachineServiceImpl) getAllScheduledModes() []inside.ScheduledMode {
	query := "SELECT * FROM machineScheduledMode"
	rows, err := res.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var modes []inside.ScheduledMode
	for rows.Next() {
		var (
			mode inside.ScheduledMode
		)

		if err := rows.Scan(&mode.Id, &mode.DeviceId, &mode.StartTime, &mode.ModeId); err != nil {
			fmt.Println("Error: ", err.Error())
			return []inside.ScheduledMode{}
		}
		modes = append(modes, mode)
	}

	return modes
}

func (res *WashingMachineServiceImpl) GetAllScheduledModesForDevice(deviceId int) []inside.ScheduledMode {
	query := "SELECT * FROM machineScheduledMode WHERE DeviceId = ?"
	rows, err := res.db.Query(query, deviceId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var modes []inside.ScheduledMode
	for rows.Next() {
		var (
			mode         inside.ScheduledMode
			startTimeStr string
		)

		if err := rows.Scan(&mode.Id, &mode.DeviceId, &startTimeStr, &mode.ModeId); err != nil {
			fmt.Println("Error: ", err.Error())
			return []inside.ScheduledMode{}
		}

		loc, _ := time.LoadLocation("Europe/Belgrade")

		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return []inside.ScheduledMode{}
		}

		currentTime := time.Now().In(loc)

		if !startTime.Before(currentTime) {
			mode.StartTime = startTimeStr
			modes = append(modes, mode)
		}
	}

	return modes

}

func (s *WashingMachineServiceImpl) generateId() int {
	id := 0
	modes := s.getAllScheduledModes()

	for _, mode := range modes {
		if mode.Id > id {
			id = mode.Id
		}
	}
	return id + 1
}

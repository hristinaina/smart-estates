package inside

import (
	"database/sql"
	"fmt"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/models/devices/inside"
	"strings"
)

type WashingMachineService interface {
	Add(estate dtos.DeviceDTO) inside.WashingMachine
	Get(id int) inside.WashingMachine
}

type WashingMachineServiceImpl struct {
	db *sql.DB
}

func NewWashingMachineService(db *sql.DB) WashingMachineService {
	return &WashingMachineServiceImpl{db: db}
}

func (s *WashingMachineServiceImpl) Get(id int) inside.WashingMachine {
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

	fmt.Println("PROSLOOOOOOOOOOO")
	// Execute the query
	rows, err := s.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return inside.WashingMachine{}
	}
	defer rows.Close()

	var wm inside.WashingMachine
	var device models.Device
	var consDevice models.ConsumptionDevice

	fmt.Println("STIGLI SMO OVDE")

	for rows.Next() {
		fmt.Println("USLI SMO OVDEEEEEEEEEEEEE")
		var startTime sql.NullString
		var name sql.NullString
		var duration sql.NullFloat64
		var temperature sql.NullFloat64
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

	return wm
}

func (s *WashingMachineServiceImpl) findModeBasedOnName(names string) []inside.Mode {
	var modes []inside.Mode

	parts := strings.Split(names, ",")

	for _, part := range parts {
		var id int
		var duration int
		var temp int
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
	return device
}

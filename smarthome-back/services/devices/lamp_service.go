package services

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"smarthome-back/dto"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
	"sort"
)

type LampService interface {
	Get(id int) (models.Lamp, error)
	GetAll() ([]models.Lamp, error)
	TurnOn(id int) (models.Lamp, error)
	TurnOff(id int) (models.Lamp, error)
	SetLightning(id int, level int) (models.Lamp, error)
	Add(dto dto.DeviceDTO) (models.Lamp, error)
	Delete(id int) (bool, error)
	GetGraphData(from, to string) ([]dtos.GraphData, error)
}

type LampServiceImpl struct {
	db         *sql.DB
	influxDb   influxdb2.Client
	repository repositories.LampRepository
}

func NewLampService(db *sql.DB, influxDb influxdb2.Client) LampService {
	return &LampServiceImpl{db: db, repository: repositories.NewLampRepository(db, influxDb)}
}

func (ls *LampServiceImpl) Get(id int) (models.Lamp, error) {
	return ls.repository.Get(id)
}

func (ls *LampServiceImpl) GetAll() ([]models.Lamp, error) {
	return ls.repository.GetAll()
}

func (ls *LampServiceImpl) TurnOn(id int) (models.Lamp, error) {
	_, err := ls.repository.UpdateIsOnState(id, true)
	if err != nil {
		return models.Lamp{}, err
	}
	lamp, err := ls.Get(id)
	return lamp, err
}

func (ls *LampServiceImpl) TurnOff(id int) (models.Lamp, error) {
	_, err := ls.repository.UpdateIsOnState(id, false)
	if err != nil {
		return models.Lamp{}, err
	}
	lamp, err := ls.Get(id)
	return lamp, err
}

func (ls *LampServiceImpl) SetLightning(id int, level int) (models.Lamp, error) {
	_, err := ls.repository.UpdateLightningState(id, level)
	if err != nil {
		return models.Lamp{}, err
	}
	lamp, err := ls.Get(id)
	return lamp, err
}

func (ls *LampServiceImpl) Add(dto dto.DeviceDTO) (models.Lamp, error) {
	device := dto.ToLamp()
	tx, err := ls.db.Begin()
	if err != nil {
		return models.Lamp{}, err
	}
	defer tx.Rollback()
	//currentTime := mysql.NullTime{
	//	Time:  time.Now(),
	//	Valid: true,
	//}

	// TODO: move transaction to repository
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.ConsumptionDevice.Device.Name, device.ConsumptionDevice.Device.Type,
		device.ConsumptionDevice.Device.RealEstate, device.ConsumptionDevice.Device.IsOnline)
	if err != nil {
		return models.Lamp{}, err
	}

	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.Lamp{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO ConsumptionDevice(DeviceId, PowerSupply, PowerConsumption)
							VALUES (?, ?, ?)`, deviceID, device.ConsumptionDevice.PowerSupply,
		device.ConsumptionDevice.PowerConsumption)
	if err != nil {
		return models.Lamp{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO Lamp(DeviceId, IsOn, LightningLevel)
							VALUES (?, ?, ?)`, deviceID, device.IsOn, device.LightningLevel)
	if err != nil {
		return models.Lamp{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.Lamp{}, err
	}
	device.ConsumptionDevice.Device.Id = int(deviceID)
	return device, nil
}

func (ls *LampServiceImpl) Delete(id int) (bool, error) {
	// checking if lamp with selected id exists
	_, err := ls.Get(id)
	if err != nil {
		return false, err
	}

	tx, err := ls.db.Begin()
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

	// TODO: move this to repository
	_, err = tx.Exec("DELETE FROM Lamp WHERE DeviceId = ?", id)
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

func (ls *LampServiceImpl) GetGraphData(from, to string) ([]dtos.GraphData, error) {
	values := make(map[float64]int)
	results := ls.repository.GetLampData(from, to)

	for results.Next() {
		if results.Record().Value() != nil {
			key := results.Record().Value().(float64)
			if isPresentInMap(values, key) {
				currentCount := getValueFromMap(values, key)
				currentCount++
				values[key] = currentCount
			} else {
				values[key] = 1
			}
		}
	}
	var graphData []dtos.GraphData
	var keys []float64
	for k := range values {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for _, k := range keys {
		g := dtos.GraphData{
			Count: values[k],
			Value: k,
		}
		graphData = append(graphData, g)
	}
	return graphData, nil
}

func isPresentInMap(mapValues map[float64]int, key float64) bool {
	if _, ok := mapValues[key]; ok {
		return true
	}
	return false
}

func getValueFromMap(mapValues map[float64]int, key float64) int {
	if value, ok := mapValues[key]; ok {
		return value
	}
	// TODO: think about this
	return -1
}

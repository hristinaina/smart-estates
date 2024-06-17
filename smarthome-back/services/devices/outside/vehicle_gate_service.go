package outside

import (
	"database/sql"
	"fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	"smarthome-back/dtos/vehicle_gate_graph"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
	"sort"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type VehicleGateService interface {
	Get(id int) (models.VehicleGate, error)
	GetAll() ([]models.VehicleGate, error)
	Open(id int) (models.VehicleGate, error)
	Close(id int) (models.VehicleGate, error)
	ToPrivate(id int) (models.VehicleGate, error)
	ToPublic(id int) (models.VehicleGate, error)
	Add(dto dtos.DeviceDTO) (models.VehicleGate, error)
	Delete(id int) (bool, error)
	GetLicensePlates(id int) ([]string, error)
	AddLicensePlate(deviceId int, licensePlate string) (string, error)
	GetAllLicensePlates() ([]string, error)
	GetLicensePlatesCount(id int, from string, filter ...string) []vehicle_gate_graph.VehicleEntriesCount
	GetLicensePlatesOutcome(id int, from string, filter ...string) map[string]int
}

type VehicleGateServiceImpl struct {
	db           *sql.DB
	influx       influxdb2.Client
	repository   repositories.VehicleGateRepository
	cacheService cache.CacheService
}

func NewVehicleGateService(db *sql.DB, influx influxdb2.Client, cacheService cache.CacheService) VehicleGateService {
	return &VehicleGateServiceImpl{db: db, influx: influx, repository: repositories.NewVehicleGateRepository(db, influx, cacheService), cacheService: cacheService}
}

func (service *VehicleGateServiceImpl) Get(id int) (models.VehicleGate, error) {
	return service.repository.Get(id)
}

func (service *VehicleGateServiceImpl) GetAll() ([]models.VehicleGate, error) {
	return service.repository.GetAll()
}

func (service *VehicleGateServiceImpl) Open(id int) (models.VehicleGate, error) {
	_, err := service.repository.UpdateIsOpen(id, true)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate, err := service.Get(id)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate.IsOpen = true
	return gate, nil
}

func (service *VehicleGateServiceImpl) Close(id int) (models.VehicleGate, error) {
	_, err := service.repository.UpdateIsOpen(id, false)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate, err := service.Get(id)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate.IsOpen = false
	return gate, nil
}

func (service *VehicleGateServiceImpl) ToPrivate(id int) (models.VehicleGate, error) {
	_, err := service.repository.UpdateMode(id, enumerations.Private)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate, err := service.Get(id)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate.Mode = enumerations.Private
	return gate, nil
}

func (service *VehicleGateServiceImpl) ToPublic(id int) (models.VehicleGate, error) {
	_, err := service.repository.UpdateMode(id, enumerations.Public)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate, err := service.Get(id)
	if err != nil {
		return models.VehicleGate{}, err
	}
	gate.Mode = enumerations.Public
	return gate, nil
}

func (service *VehicleGateServiceImpl) Add(dto dtos.DeviceDTO) (models.VehicleGate, error) {
	device := dto.ToVehicleGate()
	tx, err := service.db.Begin()
	if err != nil {
		return models.VehicleGate{}, err
	}
	defer tx.Rollback()

	// TODO: move transaction to repository
	result, err := tx.Exec(`
		INSERT INTO Device (Name, Type, RealEstate, IsOnline)
		VALUES (?, ?, ?, ?)
	`, device.ConsumptionDevice.Device.Name, device.ConsumptionDevice.Device.Type,
		device.ConsumptionDevice.Device.RealEstate, device.ConsumptionDevice.Device.IsOnline)
	if err != nil {
		return models.VehicleGate{}, err
	}

	deviceID, err := result.LastInsertId()
	if err != nil {
		return models.VehicleGate{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO ConsumptionDevice(DeviceId, PowerSupply, PowerConsumption)
							VALUES (?, ?, ?)`, deviceID, device.ConsumptionDevice.PowerSupply,
		device.ConsumptionDevice.PowerConsumption)
	if err != nil {
		return models.VehicleGate{}, err
	}

	result, err = tx.Exec(`
							INSERT INTO vehicleGate(DeviceId, IsOpen, Mode)
							VALUES (?, ?, ?)`, deviceID, device.IsOpen, device.Mode)
	if err != nil {
		return models.VehicleGate{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.VehicleGate{}, err
	}
	device.ConsumptionDevice.Device.Id = int(deviceID)

	cacheKey := fmt.Sprintf("gate_%d", device.ConsumptionDevice.Device.Id)
	if err := service.cacheService.SetToCache(cacheKey, device); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}
	err = service.cacheService.AddDevicesByRealEstate(device.ConsumptionDevice.Device.RealEstate, device.ConsumptionDevice.Device)
	return device, nil
}

func (service *VehicleGateServiceImpl) Delete(id int) (bool, error) {
	return service.repository.Delete(id)
}

func (service *VehicleGateServiceImpl) GetLicensePlates(id int) ([]string, error) {
	return service.repository.GetLicensePlates(id)
}

func (service *VehicleGateServiceImpl) AddLicensePlate(deviceId int, licensePlate string) (string, error) {
	return service.repository.AddLicensePlate(deviceId, licensePlate)
}

func (service *VehicleGateServiceImpl) GetAllLicensePlates() ([]string, error) {
	return service.repository.GetAllLicensePlates()
}

func (service *VehicleGateServiceImpl) GetLicensePlatesCount(id int, from string, filter ...string) []vehicle_gate_graph.VehicleEntriesCount {
	values := make(map[string]int)
	var result *api.QueryTableResult
	if len(filter) == 1 {
		result = service.repository.GetFromInfluxDb(id, from, filter[0])
	} else {
		result = service.repository.GetFromInfluxDb(id, from, filter[0], filter[1])
	}

	fmt.Println("Resultsssss")
	fmt.Println(result)

	for result.Next() {
		if result.Record().Value() != nil {
			key := result.Record().Value().(string)
			if service.isPresentInMap(values, key) {
				currentCount := service.getValueFromMap(values, key)
				currentCount++
				values[key] = currentCount
			} else {
				values[key] = 1
			}
		}
	}

	var graphData []vehicle_gate_graph.VehicleEntriesCount
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		g := vehicle_gate_graph.VehicleEntriesCount{
			Count:        values[k],
			LicensePlate: k,
		}
		graphData = append(graphData, g)
	}

	return graphData
}

func (service *VehicleGateServiceImpl) GetLicensePlatesOutcome(id int, from string, filter ...string) map[string]int {
	values := make(map[string]int)
	var result *api.QueryTableResult
	result = service.repository.GetFromInfluxDb(id, from, filter[0])

	fmt.Println("Resultsssss")
	fmt.Printf("Type of record: %T\n", result.Record())
	fmt.Println(result)
	values["success"] = 0
	values["not_success"] = 0
	is_success := "not_success"

	for result.Next() {
		if result.Record().Value() != nil {
			//key := result.Record().Value().(string)
			//fmt.Println(result.Record())
			successValue := result.Record().ValueByKey("Success")
			//key := result.Record().ValueByKey("Action").(string)
			if successValue == "true" {
				is_success = "success"

			} else {
				is_success = "not_success"
			}
			currentCount := service.getValueFromMap(values, is_success)
			currentCount++
			values[is_success] = currentCount
		}
	}

	return values
}

func (service *VehicleGateServiceImpl) isPresentInMap(mapValues map[string]int, key string) bool {
	if _, ok := mapValues[key]; ok {
		return true
	}
	return false
}

func (service *VehicleGateServiceImpl) getValueFromMap(mapValues map[string]int, key string) int {
	if value, ok := mapValues[key]; ok {
		return value
	}
	// TODO: think about returning -1
	return -1
}

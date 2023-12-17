package services

import (
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dto"
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
)

type VehicleGateService interface {
	Get(id int) (models.VehicleGate, error)
	GetAll() ([]models.VehicleGate, error)
	Open(id int) (models.VehicleGate, error)
	Close(id int) (models.VehicleGate, error)
	ToPrivate(id int) (models.VehicleGate, error)
	ToPublic(id int) (models.VehicleGate, error)
	Add(dto dto.DeviceDTO) (models.VehicleGate, error)
	Delete(id int) (bool, error)
	GetLicensePlates(id int) ([]string, error)
	AddLicensePlate(deviceId int, licensePlate string) (string, error)
	GetAllLicensePlates() ([]string, error)
}

type VehicleGateServiceImpl struct {
	db         *sql.DB
	influx     influxdb2.Client
	repository repositories.VehicleGateRepository
}

func NewVehicleGateService(db *sql.DB, influx influxdb2.Client) VehicleGateService {
	return &VehicleGateServiceImpl{db: db, influx: influx, repository: repositories.NewVehicleGateRepository(db, influx)}
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

func (service *VehicleGateServiceImpl) Add(dto dto.DeviceDTO) (models.VehicleGate, error) {
	device := dto.ToVehicleGate()
	tx, err := service.db.Begin()
	if err != nil {
		return models.VehicleGate{}, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println("Rollback error: ", err)
		}
	}(tx)

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
							INSERT INTO vehiclegate(DeviceId, IsOpen, Mode)
							VALUES (?, ?, ?)`, deviceID, device.IsOpen, device.Mode)
	if err != nil {
		return models.VehicleGate{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.VehicleGate{}, err
	}
	device.ConsumptionDevice.Device.Id = int(deviceID)
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

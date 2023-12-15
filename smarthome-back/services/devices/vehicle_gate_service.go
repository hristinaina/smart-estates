package services

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

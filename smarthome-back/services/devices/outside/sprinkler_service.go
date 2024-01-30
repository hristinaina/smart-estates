package outside

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
)

type SprinklerService interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() ([]models.Sprinkler, error)
	UpdateIsOn(id int, isOn bool) (bool, error)
	Delete(id int) (bool, error)
	Add(dto dtos.DeviceDTO) (models.Sprinkler, error)
}

type SprinklerServiceImpl struct {
	db         *sql.DB
	influx     influxdb2.Client
	repository repositories.SprinklerRepository
}

func NewSprinklerService(db *sql.DB, client influxdb2.Client) SprinklerService {
	return &SprinklerServiceImpl{db: db, influx: client, repository: repositories.NewSprinklerRepository(db, client)}
}

func (service *SprinklerServiceImpl) Get(id int) (models.Sprinkler, error) {
	return service.repository.Get(id)
}

func (service *SprinklerServiceImpl) GetAll() ([]models.Sprinkler, error) {
	return service.repository.GetAll()
}

func (service *SprinklerServiceImpl) UpdateIsOn(id int, isOn bool) (bool, error) {
	return service.repository.UpdateIsOn(id, isOn)
}

func (service *SprinklerServiceImpl) Delete(id int) (bool, error) {
	return service.repository.Delete(id)
}

// Add TODO: this function may not be necessary because there is implementation in Device Service for this
func (service *SprinklerServiceImpl) Add(dto dtos.DeviceDTO) (models.Sprinkler, error) {
	device := dto.ToSprinkler()
	return service.repository.Add(device)
}

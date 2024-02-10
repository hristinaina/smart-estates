package outside

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
	"strconv"
	"time"
)

type SprinklerService interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() ([]models.Sprinkler, error)
	UpdateIsOn(id int, isOn bool, email string) (bool, error)
	Delete(id int) (bool, error)
	Add(dto dtos.DeviceDTO) (models.Sprinkler, error)
	AddSpecialMode(deviceId int, dto dtos.SprinklerSpecialModeDTO) (models.SprinklerSpecialMode, error)
	GetSpecialModes(deviceId int) ([]dtos.SprinklerSpecialModeDTO, error)
	DeleteSpecialMode(id int) (bool, error)
	GetSpecialMode(id int) (models.SprinklerSpecialMode, error)
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

func (service *SprinklerServiceImpl) UpdateIsOn(id int, isOn bool, email string) (bool, error) {
	res, err := service.repository.UpdateIsOn(id, isOn)
	if err == nil {
		action := "on"
		if isOn == false {
			action = "off"
		}
		saveSprinklerToInfluxDb(service.influx, id, action, email)
	}

	return res, err
}

func (service *SprinklerServiceImpl) Delete(id int) (bool, error) {
	return service.repository.Delete(id)
}

func (service *SprinklerServiceImpl) Add(dto dtos.DeviceDTO) (models.Sprinkler, error) {
	device := dto.ToSprinkler()
	return service.repository.Add(device)
}

func (service *SprinklerServiceImpl) AddSpecialMode(deviceId int, dto dtos.SprinklerSpecialModeDTO) (models.SprinklerSpecialMode, error) {
	mode := dto.ToSprinklerSpecialMode()
	return service.repository.AddSpecialMode(deviceId, mode)
}

func (service *SprinklerServiceImpl) GetSpecialModes(deviceId int) ([]dtos.SprinklerSpecialModeDTO, error) {
	res, err := service.repository.GetSpecialModes(deviceId)
	if err != nil {
		return nil, err
	}
	var ret []dtos.SprinklerSpecialModeDTO
	for _, value := range res {
		ret = append(ret, dtos.SprinklerSpecialModeToDTO(value))
	}
	return ret, nil
}

func (service *SprinklerServiceImpl) DeleteSpecialMode(id int) (bool, error) {
	return service.repository.DeleteSpecialMode(id)
}

func (service *SprinklerServiceImpl) GetSpecialMode(deviceId int) (models.SprinklerSpecialMode, error) {
	return service.repository.GetSpecialMode(deviceId)
}

func saveSprinklerToInfluxDb(client influxdb2.Client, deviceId int, mode, user string) {
	Org := "Smart Home"
	Bucket := "bucket"
	writeAPI := client.WriteAPI(Org, Bucket)

	point := influxdb2.NewPoint("sprinkler", // table
		map[string]string{"device_id": strconv.Itoa(deviceId)}, // tag
		map[string]interface{}{"action": mode, "user_id": user},
		time.Now()) // field

	writeAPI.WritePoint(point)
	writeAPI.Flush()
}

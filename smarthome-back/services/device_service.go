package services

import (
	"database/sql"
	_ "database/sql"
	"errors"
	_ "github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dto"
	"smarthome-back/models/devices"
	"smarthome-back/mqtt_client"
	"smarthome-back/repositories"
	services "smarthome-back/services/devices"
	"strconv"
)

type DeviceService interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	Add(estate dto.DeviceDTO) (models.Device, error)
	GetAll() []models.Device
	GetConsumptionDevicesByEstateId(estateId int) ([]models.ConsumptionDevice, error)
}

type DeviceServiceImpl struct {
	db                    *sql.DB
	inflixDb              influxdb2.Client
	airConditionerService AirConditionerService
	evChargerService      EVChargerService
	homeBatteryService    HomeBatteryService
	solarPanelService     SolarPanelService
	lampService           services.LampService
	mqtt                  *mqtt_client.MQTTClient
	deviceRepository      repositories.DeviceRepository
}

func NewDeviceService(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) DeviceService {
	return &DeviceServiceImpl{db: db, airConditionerService: NewAirConditionerService(db), evChargerService: NewEVChargerService(db),
		homeBatteryService: NewHomeBatteryService(db, influxDb), lampService: services.NewLampService(db, influxDb),
		mqtt: mqtt, deviceRepository: repositories.NewDeviceRepository(db),
		solarPanelService: NewSolarPanelService(db, influxDb)}
}

func (res *DeviceServiceImpl) GetAll() []models.Device {
	return res.deviceRepository.GetAll()
}

func (res *DeviceServiceImpl) GetAllByEstateId(estateId int) []models.Device {
	return res.deviceRepository.GetAllByEstateId(estateId)
}

func (res *DeviceServiceImpl) GetConsumptionDevicesByEstateId(estateId int) ([]models.ConsumptionDevice, error) {
	return res.deviceRepository.GetConsumptionDevicesByEstateId(estateId)
}

func (res *DeviceServiceImpl) Get(id int) (models.Device, error) {
	return res.deviceRepository.Get(id)
}

func (res *DeviceServiceImpl) Add(dto dto.DeviceDTO) (models.Device, error) {
	devices, err := res.deviceRepository.GetDevicesByUserID(dto.UserId)
	if err != nil {
		return models.Device{}, err
	}
	for _, value := range devices {
		if value.Name == dto.Name {
			return models.Device{}, errors.New("Device name must be unique per user")
		}
	}
	var device models.Device
	if dto.Type == 1 {
		device = res.airConditionerService.Add(dto).ToDevice()
	} else if dto.Type == 3 {
		lamp, err := res.lampService.Add(dto)
		if err != nil {
			return models.Device{}, err
		}
		device = lamp.ToDevice()
	} else if dto.Type == 8 {
		device = res.evChargerService.Add(dto).ToDevice()
	} else if dto.Type == 7 {
		device = res.homeBatteryService.Add(dto).ToDevice()
	} else if dto.Type == 6 {
		device = res.solarPanelService.Add(dto).ToDevice()
	} else {
		device = dto.ToDevice()
		query := "INSERT INTO device (Name, Type, RealEstate, IsOnline) VALUES ( ?, ?, ?, ?);"
		result, err := res.db.Exec(query, device.Name, device.Type, device.RealEstate, device.IsOnline)
		if CheckIfError(err) {
			return models.Device{}, err
		}
		id, err := result.LastInsertId()
		device.Id = int(id)
	}

	res.mqtt.Publish(mqtt_client.TopicNewDevice+strconv.Itoa(device.Id), "new device created")
	return device, nil
}

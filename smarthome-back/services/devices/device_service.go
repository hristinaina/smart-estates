package devices

import (
	"database/sql"
	_ "database/sql"
	"errors"
	_ "fmt"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/mqtt_client"
	repositories "smarthome-back/repositories/devices"
	"smarthome-back/services"
	"smarthome-back/services/devices/energetic"
	"smarthome-back/services/devices/inside"
	"smarthome-back/services/devices/outside"
	"strconv"

	_ "github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DeviceService interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	Add(estate dtos.DeviceDTO) (models.Device, error)
	GetAll() []models.Device
	GetConsumptionDevice(id int) (models.ConsumptionDevice, error)
	GetConsumptionDevicesByEstateId(estateId int) ([]models.ConsumptionDevice, error)
	GetConsumptionDeviceDto(id int) (dtos.ConsumptionDeviceDto, error)
}

type DeviceServiceImpl struct {
	db                    *sql.DB
	inflixDb              influxdb2.Client
	airConditionerService inside.AirConditionerService
	washingMachineService inside.WashingMachineService
	evChargerService      energetic.EVChargerService
	homeBatteryService    energetic.HomeBatteryService
	solarPanelService     energetic.SolarPanelService
	ambientSensorService  inside.AmbientSensorService
	lampService           outside.LampService
	vehicleGateService    outside.VehicleGateService
	sprinklerService      outside.SprinklerService
	mqtt                  *mqtt_client.MQTTClient
	deviceRepository      repositories.DeviceRepository
	cacheService          cache.CacheService
}

func NewDeviceService(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client, cacheService cache.CacheService) DeviceService {
	return &DeviceServiceImpl{db: db, airConditionerService: inside.NewAirConditionerService(db, &cacheService), washingMachineService: inside.NewWashingMachineService(db, &cacheService), evChargerService: energetic.NewEVChargerService(db, influxDb),
		homeBatteryService: energetic.NewHomeBatteryService(db, influxDb), lampService: outside.NewLampService(db, influxDb),
		vehicleGateService: outside.NewVehicleGateService(db, influxDb), sprinklerService: outside.NewSprinklerService(db, influxDb),
		mqtt: mqtt, deviceRepository: repositories.NewDeviceRepository(db, &cacheService),
		solarPanelService: energetic.NewSolarPanelService(db, influxDb), ambientSensorService: inside.NewAmbientSensorService(db, &cacheService)}
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

func (res *DeviceServiceImpl) Add(dto dtos.DeviceDTO) (models.Device, error) {
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
	if dto.Type == 0 {
		device = res.ambientSensorService.Add(dto).ToDevice()
	} else if dto.Type == 1 {
		device = res.airConditionerService.Add(dto).ToDevice()
	} else if dto.Type == 2 { // todo uradi za ves masinu
		device = res.washingMachineService.Add(dto).ToDevice()
	} else if dto.Type == 3 {
		lamp, err := res.lampService.Add(dto)
		if err != nil {
			return models.Device{}, err
		}
		device = lamp.ToDevice()
	} else if dto.Type == 4 {
		gate, err := res.vehicleGateService.Add(dto)
		if err != nil {
			return models.Device{}, err
		}
		device = gate.ToDevice()
	} else if dto.Type == 5 {
		sprinkler, err := res.sprinklerService.Add(dto)
		if err != nil {
			return models.Device{}, err
		}
		device = sprinkler.ToDevice()
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
		if services.CheckIfError(err) {
			return models.Device{}, err
		}
		id, err := result.LastInsertId()
		device.Id = int(id)
	}

	res.mqtt.Publish(mqtt_client.TopicNewDevice+strconv.Itoa(device.Id), "new device created")
	return device, nil
}

func (res *DeviceServiceImpl) GetConsumptionDeviceDto(id int) (dtos.ConsumptionDeviceDto, error) {
	return res.deviceRepository.GetConsumptionDeviceDto(id)
}

func (res *DeviceServiceImpl) GetConsumptionDevice(id int) (models.ConsumptionDevice, error) {
	return res.deviceRepository.GetConsumptionDevice(id)
}

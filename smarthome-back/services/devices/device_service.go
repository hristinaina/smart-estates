package devices

import (
	"context"
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/dtos"
	models "smarthome-back/models/devices"
	"smarthome-back/mqtt_client"
	repositories "smarthome-back/repositories/devices"
	"smarthome-back/services"
	"smarthome-back/services/devices/energetic"
	"smarthome-back/services/devices/inside"
	"smarthome-back/services/devices/outside"
	"strconv"
	"time"
)

type DeviceService interface {
	GetAllByEstateId(id int) []models.Device
	Get(id int) (models.Device, error)
	Add(estate dtos.DeviceDTO) (models.Device, error)
	GetAll() []models.Device
	GetConsumptionDevice(id int) (models.ConsumptionDevice, error)
	GetConsumptionDevicesByEstateId(estateId int) ([]models.ConsumptionDevice, error)
	GetConsumptionDeviceDto(id int) (dtos.ConsumptionDeviceDto, error)
	GetAvailability(dto dtos.ActionGraphRequest) []time.Time
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
	mqtt                  *mqtt_client.MQTTClient
	deviceRepository      repositories.DeviceRepository
}

func NewDeviceService(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) DeviceService {
	return &DeviceServiceImpl{db: db, inflixDb: influxDb, airConditionerService: inside.NewAirConditionerService(db), washingMachineService: inside.NewWashingMachineService(db), evChargerService: energetic.NewEVChargerService(db),
		homeBatteryService: energetic.NewHomeBatteryService(db, influxDb), lampService: outside.NewLampService(db, influxDb),
		vehicleGateService: outside.NewVehicleGateService(db, influxDb),
		mqtt:               mqtt, deviceRepository: repositories.NewDeviceRepository(db),
		solarPanelService: energetic.NewSolarPanelService(db, influxDb), ambientSensorService: inside.NewAmbientSensorService(db)}
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

func (res *DeviceServiceImpl) GetAvailability(dto dtos.ActionGraphRequest) []time.Time {
	influxOrg := "Smart Home"
	influxBucket := "bucket"

	// Create InfluxDB query API
	queryAPI := res.inflixDb.QueryAPI(influxOrg)
	// Define your InfluxDB query with conditions
	query := fmt.Sprintf(
		`from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r["_measurement"] == "device_status" and r["device_id"] == "%s" and r["_field"] == "status" and r["_value"] == 1)
		  |> map(fn: (r) => ({
			time: r["_time"]
		  }))`,
		influxBucket, dto.StartDate, dto.EndDate, strconv.Itoa(dto.DeviceId),
	)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Error executing InfluxDB query: %v\n", err)
		return nil
	}

	defer result.Close()
	var times []time.Time
	fmt.Println("Printing influxDB data...")
	for result.Next() {
		fmt.Println("------------------------")
		fmt.Println(result.Record())
		fmt.Println(result.Record().Time())
		times = append(times, result.Record().Time())
		if err := result.Err(); err != nil {
			fmt.Println("ERROR happened")
		}
	}
	return times
}

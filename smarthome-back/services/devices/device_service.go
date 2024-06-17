package devices

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	"strings"
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
	GetAvailability(dto dtos.ActionGraphRequest) map[time.Time]float64
}

type DeviceServiceImpl struct {
	db                    *sql.DB
	influxDb              influxdb2.Client
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
	return &DeviceServiceImpl{db: db, airConditionerService: inside.NewAirConditionerService(db, &cacheService), washingMachineService: inside.NewWashingMachineService(db, &cacheService), evChargerService: energetic.NewEVChargerService(db, influxDb, cacheService),
		homeBatteryService: energetic.NewHomeBatteryService(db, influxDb, cacheService), lampService: outside.NewLampService(db, influxDb, cacheService),
		vehicleGateService: outside.NewVehicleGateService(db, influxDb, cacheService), sprinklerService: outside.NewSprinklerService(db, influxDb, cacheService),
		mqtt: mqtt, deviceRepository: repositories.NewDeviceRepository(db, influxDb, &cacheService),
		solarPanelService: energetic.NewSolarPanelService(db, influxDb, cacheService), ambientSensorService: inside.NewAmbientSensorService(db, &cacheService)}
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

func (res *DeviceServiceImpl) GetAvailability(dto dtos.ActionGraphRequest) map[time.Time]float64 {
	onlineTimes := res.deviceRepository.GetAvailability(dto, "1")
	offlineTimes := res.deviceRepository.GetAvailability(dto, "0")

	if dto.EndDate == "-1" {
		return res.getAvailabilityPerHour(onlineTimes, offlineTimes)
	} else {
		if GetNumOfPassedDays(ParseDate(dto.StartDate), ParseDate(dto.EndDate)) >= 2 {
			return res.getAvailabilityPerDay(onlineTimes, offlineTimes)
		} else {
			return res.getAvailabilityPerHour(onlineTimes, offlineTimes)
		}
	}
}

func (res *DeviceServiceImpl) GetTotalOnlineOfflineHours(dto dtos.ActionGraphRequest) (float64, float64) {
	onlineTimes := res.deviceRepository.GetAvailability(dto, "1")
	offlineTimes := res.deviceRepository.GetAvailability(dto, "0")

	totalOnlineDuration := res.getTotalOnlineDuration(onlineTimes, offlineTimes).Hours()
	if dto.EndDate == "-1" {
		startDate := strings.TrimRight(dto.StartDate, "h")
		passedTime, err := strconv.Atoi(startDate)
		if err != nil {
			fmt.Printf("Error parsing duration: %v\n", err)
			return 0, 0
		}
		return totalOnlineDuration, float64(passedTime)*(-1) - totalOnlineDuration
	} else {
		totalHours := GetNumOfPassedHours(ParseDate(dto.StartDate), ParseDate(dto.EndDate))
		return totalOnlineDuration, totalHours - totalOnlineDuration
	}
}

func (res *DeviceServiceImpl) getAvailabilityPerHour(online, offline []time.Time) map[time.Time]float64 {
	totalDurationPerHour := make(map[time.Time]float64)

	length := Min(len(online), len(offline))
	for i := 0; i < length; i++ {
		startHour := online[i].Truncate(time.Hour)
		stopHour := offline[i].Truncate(time.Hour)

		flag := false
		if (online[i].Day() != offline[i].Day()) || (online[i].Month() != offline[i].Month()) {
			flag = true
		}

		if (online[i].Hour() != offline[i].Hour()) || flag {
			for currentHour := startHour; currentHour.Before(stopHour); currentHour = currentHour.Add(time.Hour) {
				nextHour := currentHour.Add(time.Hour)
				if nextHour.After(stopHour) {
					nextHour = stopHour // this is why 'if totalDurationPerHour[...] > 1' ... is needed
				}

				durationOnline := nextHour.Sub(online[i])
				hoursOnline := durationOnline.Hours()
				totalDurationPerHour[currentHour] += hoursOnline
				if totalDurationPerHour[currentHour] > 1 {
					totalDurationPerHour[currentHour] = 1
				}
			}
		} else {
			startTime := max(online[i], offline[i])
			endTime := min(offline[i], online[i])
			durationOnline := max(endTime, startTime).Sub(min(online[i], offline[i]))
			totalDurationPerHour[startHour] += durationOnline.Hours()
		}
	}

	//for hour, totalOnlineTime := range totalDurationPerHour {
	//	fmt.Printf("Hour: %v, Total Online Time: %.2f\n", hour.Format("2006-01-02 15:04:05"), totalOnlineTime)
	//}

	return totalDurationPerHour
}

func (res *DeviceServiceImpl) getAvailabilityPerDay(online, offline []time.Time) map[time.Time]float64 {
	totalDurationPerDay := make(map[time.Time]float64)

	length := Min(len(online), len(offline))
	offlineInd := 0
	for i := 0; i < length; i++ {
		if offlineInd >= length {
			break
		}
		startDate := online[i].Truncate(24 * time.Hour)
		stopDate := offline[offlineInd].Truncate(24 * time.Hour)
		if offline[offlineInd].Before(online[i]) {
			for {
				offlineInd++
				if offlineInd >= length {
					break
				}
				stopDate = offline[offlineInd].Truncate(24 * time.Hour)
				if offline[offlineInd].After(online[i]) {
					break
				}
			}

		}

		if stopDate.After(startDate) {
			if GetNumOfPassedDays(online[i], offline[offlineInd]) > 1 {
				// TODO: won't work for every case (offline more than 1 day after)
				// problem is our data in influx is not right
				continue
			}
			// total time for first date
			firstDayDuration := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59,
				0, startDate.Location()).Sub(online[i])
			totalDurationPerDay[startDate] += firstDayDuration.Hours()

			// total time for second date
			secondDayDuration := offline[offlineInd].Sub(time.Date(stopDate.Year(), stopDate.Month(), stopDate.Day(), 0,
				0, 0, 0, stopDate.Location()))
			totalDurationPerDay[stopDate] += secondDayDuration.Hours()
		} else {
			// if online and offline are the same date
			durationOnline := offline[offlineInd].Sub(online[i])
			totalDurationPerDay[startDate] += durationOnline.Hours()
		}
		offlineInd++
		//fmt.Printf("Index: %d, Online: %v, Offline: %v, StartDate: %v, StopDate: %v, Duration: %v\n", i, online[i], offline[i], startDate, stopDate, totalDurationPerDay[startDate])
	}
	for date, totalOnlineTime := range totalDurationPerDay {
		fmt.Printf("Date: %v, Total Online Time: %v\n", date.Format("2006-01-02"), totalOnlineTime)
	}
	return totalDurationPerDay
}

func (res *DeviceServiceImpl) getTotalOnlineDuration(onlineTimes, offlineTimes []time.Time) time.Duration {
	total := time.Duration(0)
	length := Min(len(onlineTimes), len(offlineTimes))

	for i := 0; i < length; i++ {
		duration := offlineTimes[i].Sub(onlineTimes[i])
		total += duration
	}

	return total
}

func Min(firstNum, secondNum int) int {
	if firstNum <= secondNum {
		return firstNum
	}
	return secondNum
}

func GetNumOfPassedDays(start, end time.Time) float64 {
	return end.Sub(start).Hours() / 24.0
}

func GetNumOfPassedHours(start, end time.Time) float64 {
	return end.Sub(start).Hours()
}

func ParseDate(date string) time.Time {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Errorf("error happened %s", err)
		return time.Time{}
	}
	return parsedDate
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// Utility function to find the minimum of two times
func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

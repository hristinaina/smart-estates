package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"smarthome-back/cache"
	"smarthome-back/enumerations"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type VehicleGateRepository interface {
	Get(id int) (models.VehicleGate, error)
	GetAll() ([]models.VehicleGate, error)
	UpdateIsOpen(id int, isOpen bool) (bool, error)
	UpdateMode(id int, mode enumerations.VehicleGateMode) (bool, error)
	Delete(id int) (bool, error)
	GetLicensePlates(id int) ([]string, error)
	AddLicensePlate(id int, licensePlate string) (string, error)
	GetAllLicensePlates() ([]string, error)
	PostNewVehicleGateValue(gate models.VehicleGate, action string, success bool, licensePlate string)
	GetFromInfluxDb(id int, from string, filter ...string) *api.QueryTableResult
}

type VehicleGateRepositoryImpl struct {
	db           *sql.DB
	influx       influxdb2.Client
	cacheService cache.CacheService
}

func NewVehicleGateRepository(db *sql.DB, influx influxdb2.Client, cacheService cache.CacheService) VehicleGateRepository {
	return &VehicleGateRepositoryImpl{db: db, influx: influx, cacheService: cacheService}
}

func (repo *VehicleGateRepositoryImpl) Get(id int) (models.VehicleGate, error) {
	cacheKey := fmt.Sprintf("gate_%d", id)

	var gate models.VehicleGate
	if found, err := repo.cacheService.GetFromCache(cacheKey, &gate); found {
		return gate, err
	}

	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, v.IsOpen, v.Mode
			  FROM vehicleGate v 
    		  JOIN ConsumptionDevice ON v.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id
   			  WHERE Device.Id = ? `

	rows, err := repo.db.Query(query, id)
	if repositories.IsError(err) {
		return models.VehicleGate{}, err
	}
	defer rows.Close()

	gates, err := repo.scanRows(rows)
	if repositories.IsError(err) {
		return models.VehicleGate{}, err
	}
	gate = gates[0]
	licensePlates, err := repo.GetLicensePlates(gate.ConsumptionDevice.Device.Id)
	if repositories.CheckIfError(err) {
		return models.VehicleGate{}, err
	}
	gate.LicensePlates = licensePlates

	if err := repo.cacheService.SetToCache(cacheKey, gate); err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Saved data in cache.")
	}

	return gate, nil
}

func (repo *VehicleGateRepositoryImpl) GetAll() ([]models.VehicleGate, error) {
	query := `SELECT Device.Id, Device.Name, Device.Type, Device.RealEstate, Device.IsOnline,
       		  ConsumptionDevice.PowerSupply, ConsumptionDevice.PowerConsumption, v.IsOpen, v.Mode
			  FROM vehiclegate v 
    		  JOIN ConsumptionDevice ON v.DeviceId = ConsumptionDevice.DeviceId
   			  JOIN Device ON ConsumptionDevice.DeviceId = Device.Id`
	rows, err := repo.db.Query(query)
	if repositories.IsError(err) {
		return nil, err
	}
	defer rows.Close()

	gates, err := repo.scanRows(rows)
	gatesWithPlates := make([]models.VehicleGate, 0)
	if err != nil {
		return nil, err
	}
	for _, gate := range gates {
		licensePlates, err := repo.GetLicensePlates(gate.ConsumptionDevice.Device.Id)
		if repositories.CheckIfError(err) {
			return nil, err
		}
		gate.LicensePlates = licensePlates
		gatesWithPlates = append(gatesWithPlates, gate)
	}
	return gatesWithPlates, nil
}

func (repo *VehicleGateRepositoryImpl) UpdateIsOpen(id int, isOpen bool) (bool, error) {
	query := `UPDATE VehicleGate v
			  JOIN ConsumptionDevice cd ON v.DeviceId = cd.DeviceId
			  JOIN Device d ON cd.DeviceId = d.Id
			  SET v.IsOpen = ?
			  WHERE d.Id = ?`
	_, err := repo.db.Query(query, isOpen, id)
	if repositories.IsError(err) {
		return false, err
	}
	return true, nil
}

func (repo *VehicleGateRepositoryImpl) UpdateMode(id int, mode enumerations.VehicleGateMode) (bool, error) {
	queryMode := 0
	if mode == enumerations.Public {
		queryMode = 1
	}
	query := `UPDATE VehicleGate v
			  JOIN ConsumptionDevice cd ON v.DeviceId = cd.DeviceId
			  JOIN Device d ON cd.DeviceId = d.Id
			  SET v.Mode = ?
			  WHERE d.Id = ?`

	_, err := repo.db.Exec(query, queryMode, id)
	if repositories.IsError(err) {
		return false, err
	}
	return true, nil
}

func (repo *VehicleGateRepositoryImpl) Delete(id int) (bool, error) {
	_, err := repo.Get(id)
	if err != nil {
		return false, err
	}

	tx, err := repo.db.Begin()
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

	_, err = tx.Exec("DELETE FROM licenseplate WHERE DeviceId = ?", id)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("DELETE FROM VehicleGate WHERE DeviceId = ?", id)
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

func (repo *VehicleGateRepositoryImpl) GetLicensePlates(id int) ([]string, error) {
	query := `SELECT PlateNumber FROM licensePlate WHERE DeviceId = ?`

	rows, err := repo.db.Query(query, id)
	if repositories.IsError(err) {
		return nil, err
	}
	defer rows.Close()

	licensePlates, err := repo.scanLicensePlateRows(rows)
	return licensePlates, nil
}

func (repo *VehicleGateRepositoryImpl) AddLicensePlate(id int, licensePlate string) (string, error) {
	query := `INSERT INTO licensePlate (DeviceId, PlateNumber) VALUES (?, ?)`
	_, err := repo.db.Exec(query, id, licensePlate)
	if repositories.CheckIfError(err) {
		return "", err
	}
	return licensePlate, nil
}

func (repo *VehicleGateRepositoryImpl) GetAllLicensePlates() ([]string, error) {
	query := `SELECT DISTINCT licensePlate.PlateNumber FROM licensePlate`

	rows, err := repo.db.Query(query)
	if repositories.IsError(err) {
		return make([]string, 0), err
	}
	defer rows.Close()

	licensePlates, err := repo.scanLicensePlateRows(rows)
	return licensePlates, nil
}

func (repo *VehicleGateRepositoryImpl) PostNewVehicleGateValue(gate models.VehicleGate, action string, success bool,
	licensePlate string) {
	client := repo.influx
	writeAPI := client.WriteAPIBlocking("Smart Home", "bucket")
	tags := map[string]string{
		"Id":      strconv.Itoa(gate.ConsumptionDevice.Device.Id),
		"Action":  action,
		"Success": strconv.FormatBool(success),
	}
	fields := map[string]interface{}{
		"LicensePlate": licensePlate,
	}
	point := write.NewPoint("gates", tags, fields, time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}
}

func (repo *VehicleGateRepositoryImpl) GetFromInfluxDb(id int, from string, filter ...string) *api.QueryTableResult {
	client := repo.influx
	queryAPI := client.QueryAPI("Smart Home")
	query := ""
	queryId := strconv.Itoa(id)
	if len(filter) == 1 {
		query = fmt.Sprintf(`from(bucket: "bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "gates" and r.Id == "%s")`, from, filter[0], queryId)
	} else {
		query = fmt.Sprintf(`from(bucket: "bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "gates" and r.Id == "%s" and r._field == "LicensePlate"
			and r._value == "%s")`,
			from, filter[0], queryId, filter[1])
	}

	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Printing influxDB data...")
	//for results.Next() {
	//	fmt.Println("------------------------")
	//	fmt.Println(results.Record())
	//}
	//if err := results.Err(); err != nil {
	//	log.Fatal(err)
	//}

	return results
}

func (repo *VehicleGateRepositoryImpl) scanLicensePlateRows(rows *sql.Rows) ([]string, error) {
	licensePlates := make([]string, 0)
	for rows.Next() {
		var (
			licensePlate string
		)
		if err := rows.Scan(&licensePlate); err != nil {
			fmt.Println("Error: ", err.Error())
			return nil, err
		}
		licensePlates = append(licensePlates, licensePlate)
	}

	return licensePlates, nil
}

// scanRows parses value from db to desired model - in this case to vehicle gate
func (repo *VehicleGateRepositoryImpl) scanRows(rows *sql.Rows) ([]models.VehicleGate, error) {
	var gates []models.VehicleGate
	for rows.Next() {
		var (
			device     models2.Device
			consDevice models2.ConsumptionDevice
			gate       models.VehicleGate
		)

		if err := rows.Scan(&device.Id, &device.Name, &device.Type, &device.RealEstate,
			&device.IsOnline, &consDevice.PowerSupply, &consDevice.PowerConsumption, &gate.IsOpen, &gate.Mode); err != nil {
			fmt.Println("Error: ", err.Error())
			return []models.VehicleGate{}, err
		}
		consDevice.Device = device
		gate.ConsumptionDevice = consDevice
		gates = append(gates, gate)
	}

	return gates, nil
}

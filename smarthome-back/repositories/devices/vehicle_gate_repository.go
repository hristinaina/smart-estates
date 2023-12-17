package repositories

import (
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"smarthome-back/enumerations"
	models2 "smarthome-back/models/devices"
	models "smarthome-back/models/devices/outside"
	"smarthome-back/repositories"
)

type VehicleGateRepository interface {
	Get(id int) (models.VehicleGate, error)
	GetAll() ([]models.VehicleGate, error)
	UpdateIsOpen(id int, isOpen bool) (bool, error)
	UpdateMode(id int, mode enumerations.VehicleGateMode) (bool, error)
	Delete(id int) (bool, error)
	GetLicensePlates(id int) ([]string, error)
	AddLicensePlate(id int, licensePlate string) (string, error)
}

type VehicleGateRepositoryImpl struct {
	db     *sql.DB
	influx influxdb2.Client
}

func NewVehicleGateRepository(db *sql.DB, influx influxdb2.Client) VehicleGateRepository {
	return &VehicleGateRepositoryImpl{db: db, influx: influx}
}

func (repo *VehicleGateRepositoryImpl) Get(id int) (models.VehicleGate, error) {
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
	gate := gates[0]
	licensePlates, err := repo.GetLicensePlates(gate.ConsumptionDevice.Device.Id)
	if repositories.CheckIfError(err) {
		return models.VehicleGate{}, err
	}
	gate.LicensePlates = licensePlates
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
	return licensePlates, err
}

func (repo *VehicleGateRepositoryImpl) AddLicensePlate(id int, licensePlate string) (string, error) {
	query := `INSERT INTO licensePlate (DeviceId, PlateNumber) VALUES (?, ?)`
	_, err := repo.db.Exec(query, id, licensePlate)
	if repositories.CheckIfError(err) {
		return "", err
	}
	return licensePlate, err
}

func (repo *VehicleGateRepositoryImpl) scanLicensePlateRows(rows *sql.Rows) ([]string, error) {
	licensePlates := make([]string, 0)
	for rows.Next() {
		var (
			//id           int
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

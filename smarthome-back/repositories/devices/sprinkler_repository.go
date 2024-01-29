package repositories

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	models "smarthome-back/models/devices/outside"
)

type SprinklerRepository interface {
	Get(id int) (models.Sprinkler, error)
	GetAll() (models.Sprinkler, error)
	UpdateIsOn(isOn bool) (bool, error)
	Delete(id int) (bool, error)
	AddSpecialMode(mode models.SprinklerSpecialMode) (models.Sprinkler, error)
}

type SprinklerRepositoryImpl struct {
	db     *sql.DB
	influx influxdb2.Client
}

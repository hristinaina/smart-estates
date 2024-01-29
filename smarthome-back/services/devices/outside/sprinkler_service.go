package outside

import (
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	models "smarthome-back/models/devices/outside"
	repositories "smarthome-back/repositories/devices"
)

type SprinklerService interface {
	Get(id int) (models.Sprinkler, error)
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

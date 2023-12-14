package services

import (
	"database/sql"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AmbientSensorService interface {
}

type AmbientSensorServiceImpl struct {
	db *sql.DB
}

func NewAmbientSensorService(db *sql.DB) AmbientSensorService {
	return &AmbientSensorServiceImpl{db: db}
}

func (as *AmbientSensorServiceImpl) GetValues(client mqtt.Client) {

}

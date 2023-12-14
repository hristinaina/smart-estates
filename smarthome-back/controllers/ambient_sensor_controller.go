package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/mqtt_client"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
)

type AmbientSensorController struct {
	service services.AmbientSensorService
	mqtt    *mqtt_client.MQTTClient
}

func NewAmbientSensorController(db *sql.DB, mqtt *mqtt_client.MQTTClient) AmbientSensorController {
	return AmbientSensorController{service: services.NewAmbientSensorService(db), mqtt: mqtt}
}

type AmbientSensor struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func (uc AmbientSensorController) GetValueForHour(c *gin.Context) {
	results := mqtt_client.GetLastOneHourValues(uc.mqtt.GetInflux(), c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"result": results})
}

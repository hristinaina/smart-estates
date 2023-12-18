package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/mqtt_client"
	"smarthome-back/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AirConditionerController struct {
	service services.AirConditionerService
	mqtt    *mqtt_client.MQTTClient
}

func NewAirConditionerController(db *sql.DB, mqtt *mqtt_client.MQTTClient) AirConditionerController {
	return AirConditionerController{service: services.NewAirConditionerService(db), mqtt: mqtt}
}

func (uc AirConditionerController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (ac AirConditionerController) GetHistoryData(c *gin.Context) {
	results := mqtt_client.QueryDeviceData(ac.mqtt.GetInflux(), c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"result": results})
}

package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/dto"
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
	var data dto.ActionGraphRequest
	// convert json object to model device
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	results := mqtt_client.QueryDeviceData(ac.mqtt.GetInflux(), data)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

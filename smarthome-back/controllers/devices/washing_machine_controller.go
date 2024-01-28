package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/controllers"
	"smarthome-back/mqtt_client"
	"smarthome-back/services/devices/inside"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WashingMachineController struct {
	service inside.WashingMachineService
	mqtt    *mqtt_client.MQTTClient
}

func NewWashingMachineController(db *sql.DB, mqtt *mqtt_client.MQTTClient) WashingMachineController {
	return WashingMachineController{service: inside.NewWashingMachineService(db), mqtt: mqtt}
}

func (uc WashingMachineController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

// request body
type ScheduledModeDTO struct {
	DeviceId  int
	StartTime string
	ModeId    int
}

func (uc WashingMachineController) AddScheduledMode(c *gin.Context) {
	var input ScheduledModeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	err := uc.service.AddScheduledMode(input.DeviceId, input.ModeId, input.StartTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "You have successfully scheduled the mode"})
}

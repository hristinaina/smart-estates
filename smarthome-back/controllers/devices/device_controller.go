package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/controllers"
	"smarthome-back/dtos"
	"smarthome-back/mqtt_client"
	"smarthome-back/services/devices"
	"strconv"
)

type DeviceController struct {
	service devices.DeviceService
}

func NewDeviceController(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) DeviceController {
	return DeviceController{
		service: devices.NewDeviceService(db, mqtt, influxDb)}
}

func (dc DeviceController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device, err := dc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (dc DeviceController) GetAll(c *gin.Context) {
	devices := dc.service.GetAll()
	if devices == nil {
		fmt.Println("No devices found")
		c.JSON(http.StatusBadRequest, "No devices found")
	}
	c.JSON(http.StatusOK, devices)
}

func (dc DeviceController) GetAllByEstateId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("estateId"))
	controllers.CheckIfError(err, c)
	devices := dc.service.GetAllByEstateId(id)
	c.JSON(http.StatusOK, devices)
}

func (dc DeviceController) Add(c *gin.Context) {
	var deviceDTO dtos.DeviceDTO
	// convert json object to model device
	if err := c.BindJSON(&deviceDTO); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	device, err := dc.service.Add(deviceDTO)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, device)
	}
}

func (dc DeviceController) GetConsumptionDeviceDto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	dto, err := dc.service.GetConsumptionDeviceDto(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, dto)
}

func (dc DeviceController) GetConsumptionDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	dto, err := dc.service.GetConsumptionDevice(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, dto)
}

func (dc DeviceController) GetAvailability(c *gin.Context) {
	var deviceDTO dtos.ActionGraphRequest
	// convert json object to model device
	if err := c.BindJSON(&deviceDTO); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	dc.service.GetAvailability(deviceDTO)

	c.JSON(200, "ok")
}

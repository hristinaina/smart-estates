package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	"smarthome-back/dtos"
	"smarthome-back/mqtt_client"
	"smarthome-back/services/devices"
	"strconv"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DeviceController struct {
	service devices.DeviceService
}

func NewDeviceController(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client, cacheService cache.CacheService) DeviceController {
	return DeviceController{
		service: devices.NewDeviceService(db, mqtt, influxDb, cacheService)}
}

func (uc DeviceController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device, err := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (uc DeviceController) GetAll(c *gin.Context) {
	devices := uc.service.GetAll()
	if devices == nil {
		fmt.Println("No devices found")
		c.JSON(http.StatusBadRequest, "No devices found")
	}
	c.JSON(http.StatusOK, devices)
}

func (rec DeviceController) GetAllByEstateId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("estateId"))
	controllers.CheckIfError(err, c)
	devices := rec.service.GetAllByEstateId(id)
	c.JSON(http.StatusOK, devices)
}

func (rec DeviceController) Add(c *gin.Context) {
	var deviceDTO dtos.DeviceDTO
	// convert json object to model device
	if err := c.BindJSON(&deviceDTO); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	device, err := rec.service.Add(deviceDTO)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, device)
	}
}

func (rec DeviceController) GetConsumptionDeviceDto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	dto, err := rec.service.GetConsumptionDeviceDto(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, dto)
}

func (rec DeviceController) GetConsumptionDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	dto, err := rec.service.GetConsumptionDevice(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, dto)
}

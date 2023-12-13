package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/dto"
	"smarthome-back/mqtt_client"
	"smarthome-back/services"
	"strconv"
)

type DeviceController struct {
	service services.DeviceService
}

func NewDeviceController(db *sql.DB, mqtt *mqtt_client.MQTTClient, influxDb influxdb2.Client) DeviceController {
	return DeviceController{
		service: services.NewDeviceService(db, mqtt, influxDb)}
}

func (uc DeviceController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
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
	CheckIfError(err, c)
	devices := rec.service.GetAllByEstateId(id)
	c.JSON(http.StatusOK, devices)
}

func (rec DeviceController) Add(c *gin.Context) {
	var deviceDTO dto.DeviceDTO
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

package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/models/devices"
	"smarthome-back/services"
	"strconv"
)

type DeviceController struct {
	service services.DeviceService
}

func NewDeviceController(db *sql.DB) DeviceController {
	return DeviceController{service: services.NewDeviceService(db)}
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

func (rec DeviceController) GetAllByEstateId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("estateId"))
	CheckIfError(err, c)
	realEstates := rec.service.GetAllByEstateId(id)
	c.JSON(http.StatusOK, realEstates)
}

func (rec DeviceController) Add(c *gin.Context) {
	var device models.Device

	// convert json object to model device
	if err := c.BindJSON(&device); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	device = rec.service.Add(device)
	c.JSON(http.StatusOK, device)
}

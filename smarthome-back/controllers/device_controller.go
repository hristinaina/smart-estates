package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/dto"
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

func (uc DeviceController) GetAll(c *gin.Context) {
	devices := uc.service.GetAll()
	if devices == nil {
		fmt.Println("Error happened!")
		c.JSON(http.StatusBadRequest, "Error happened!")
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
	device := rec.service.Add(deviceDTO)
	c.JSON(http.StatusOK, device)
}

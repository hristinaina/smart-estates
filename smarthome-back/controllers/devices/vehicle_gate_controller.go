package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/controllers"
	dto2 "smarthome-back/dto"
	"smarthome-back/dtos"
	services "smarthome-back/services/devices"
	"strconv"
)

type VehicleGateController struct {
	service services.VehicleGateService
}

func NewVehicleGateController(db *sql.DB, influx influxdb2.Client) VehicleGateController {
	return VehicleGateController{service: services.NewVehicleGateService(db, influx)}
}

func (controller VehicleGateController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	gate, err := controller.service.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(200, gate)
}

func (controller VehicleGateController) GetAll(c *gin.Context) {
	gates, err := controller.service.GetAll()
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, gates)
}

func (controller VehicleGateController) Open(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	gate, err := controller.service.Open(id)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, gate)
}

func (controller VehicleGateController) Close(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	gate, err := controller.service.Close(id)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, gate)
}

func (controller VehicleGateController) ToPrivate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	gate, err := controller.service.ToPrivate(id)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, gate)
}

func (controller VehicleGateController) ToPublic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	gate, err := controller.service.ToPublic(id)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, gate)
}

func (controller VehicleGateController) Add(c *gin.Context) {
	var dto dto2.DeviceDTO

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gate, err := controller.service.Add(dto)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, gate)
}

func (controller VehicleGateController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	isDeleted, err := controller.service.Delete(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Vehicle gate with selected id could not be found."})
		return
	}
	c.JSON(204, isDeleted)
}

func (controller VehicleGateController) GetLicensePlates(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	licensePlates, err := controller.service.GetLicensePlates(id)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, licensePlates)
}

func (controller VehicleGateController) AddLicensePlate(c *gin.Context) {
	var dto dtos.LicensePlate

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	licensePlate, err := controller.service.AddLicensePlate(dto.DeviceId, dto.LicensePlate)
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, licensePlate)
}

func (controller VehicleGateController) GetAllLicensePlates(c *gin.Context) {
	licensePlates, err := controller.service.GetAllLicensePlates()
	if controllers.CheckIfError(err, c) {
		return
	}
	c.JSON(200, licensePlates)
}

package controllers

import (
	"database/sql"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	"smarthome-back/dtos"
	"smarthome-back/dtos/vehicle_gate_graph"
	services "smarthome-back/services/devices/outside"
	"strconv"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type VehicleGateController struct {
	service services.VehicleGateService
}

func NewVehicleGateController(db *sql.DB, influx influxdb2.Client, cacheService cache.CacheService) VehicleGateController {
	return VehicleGateController{service: services.NewVehicleGateService(db, influx, cacheService)}
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
	var dto dtos.DeviceDTO

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

func (controller VehicleGateController) GetLicencePlatesCount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	from := c.Param("from")
	to := c.Param("to")
	licensePlate := c.Param("license-plate")

	var result []vehicle_gate_graph.VehicleEntriesCount
	if licensePlate == "-1" {
		result = controller.service.GetLicensePlatesCount(id, from, to)
	} else {
		result = controller.service.GetLicensePlatesCount(id, from, to, licensePlate)
	}
	c.JSON(200, result)
}

func (controller VehicleGateController) GetEntriesOutcome(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	from := c.Param("from")
	to := c.Param("to")

	var result map[string]int
	result = controller.service.GetLicensePlatesOutcome(id, from, to)

	c.JSON(200, result)
}

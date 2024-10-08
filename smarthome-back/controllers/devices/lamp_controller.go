package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	dto "smarthome-back/dtos"
	_ "smarthome-back/models/devices/outside"
	services "smarthome-back/services/devices/outside"
	"strconv"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type LampController struct {
	service services.LampService
}

func NewLampController(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) LampController {
	return LampController{service: services.NewLampService(db, influxDb, cacheService)}
}

func (lc LampController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) GetAll(c *gin.Context) {
	lamps, err := lc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamps)
}

func (lc LampController) TurnOn(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.TurnOn(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) TurnOff(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.TurnOff(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) SetLightning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	level, err := strconv.Atoi(c.Param("level"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.SetLightning(id, level)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) Add(c *gin.Context) {
	var dto dto.DeviceDTO

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	lamp, err := lc.service.Add(dto)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	isDeleted, err := lc.service.Delete(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Lamp with selected id not found"})
		return
	}

	c.JSON(204, isDeleted)
}

func (lc LampController) GetGraphData(c *gin.Context) {
	cId := c.Param("id")
	id, err := strconv.Atoi(cId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Parsing string to int created an error"})
		return
	}
	//var dateRange dtos.DateRange
	//if err = c.BindJSON(&dateRange); err != nil {
	//	c.JSON(400, gin.H{"error": "Invalid JSON"})
	//	return
	//}

	from := c.Param("from")
	to := c.Param("to")
	if to == "-1" {
		to = ""
	}

	data, err := lc.service.GetGraphData(id, from, to)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, gin.H{"data": data})
}

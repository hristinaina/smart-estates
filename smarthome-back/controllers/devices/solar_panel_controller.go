package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	"smarthome-back/dtos"
	"smarthome-back/services/devices/energetic"
	"strconv"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SolarPanelController struct {
	service energetic.SolarPanelService
}

func NewSolarPanelController(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) SolarPanelController {
	return SolarPanelController{service: energetic.NewSolarPanelService(db, influxDb, cacheService)}
}

func (uc SolarPanelController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (uc SolarPanelController) GetGraphData(c *gin.Context) {
	var data dtos.ActionGraphRequest
	// convert json object to model device
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	graphData, err := uc.service.GetGraphData(data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}
	c.JSON(http.StatusOK, graphData)
}

func (uc SolarPanelController) GetValueFromLastMinute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	graphData, err := uc.service.GetValueFromLastMinute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}
	c.JSON(http.StatusOK, graphData)
}

func (uc SolarPanelController) GetProductionForSP(c *gin.Context) {
	var data dtos.ActionGraphRequest
	// convert json object to model device
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	graphData := uc.service.GetProductionForSP(data)

	c.JSON(http.StatusOK, graphData)
}

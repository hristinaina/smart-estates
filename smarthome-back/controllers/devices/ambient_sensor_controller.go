package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	"smarthome-back/mqtt_client"
	"smarthome-back/services/devices/inside"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type AmbientSensorController struct {
	service  inside.AmbientSensorService
	influxDb influxdb2.Client
}

func NewAmbientSensorController(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) AmbientSensorController {
	return AmbientSensorController{service: inside.NewAmbientSensorService(db, &cacheService), influxDb: influxDb}
}

func (as AmbientSensorController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device, err := as.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

type AmbientSensor struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func (uc AmbientSensorController) GetValueForHour(c *gin.Context) {
	results := mqtt_client.GetLastOneHourValues(uc.influxDb, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"result": results})
}

type TimeInput struct {
	Time string `json:"time" binding:"required"`
}

func (uc AmbientSensorController) GetValueForSelectedTime(c *gin.Context) {
	var input TimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	results := mqtt_client.GetValuesForSelectedTime(uc.influxDb, input.Time, c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"result": results})
}

type DateInput struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func (uc AmbientSensorController) GetValuesForDate(c *gin.Context) {
	var input DateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.Start)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}

	endDate, err := time.Parse("2006-01-02", input.End)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}

	results := mqtt_client.GetValuesForDate(uc.influxDb, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"result": results})
}

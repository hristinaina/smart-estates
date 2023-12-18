package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/services"
	"strconv"
	"time"
)

type HomeBatteryController struct {
	service services.HomeBatteryService
}

func NewHomeBatteryController(db *sql.DB, influxDb influxdb2.Client) HomeBatteryController {
	return HomeBatteryController{service: services.NewHomeBatteryService(db, influxDb)}
}

func (uc HomeBatteryController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (uc HomeBatteryController) GetConsumptionForLastHour(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	results := uc.service.GetConsumptionForLastHour(id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc HomeBatteryController) GetConsumptionForSelectedTime(c *gin.Context) {
	var input TimeInput
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetConsumptionForSelectedTime(input.Time, id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc HomeBatteryController) GetConsumptionForSelectedDate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
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
	results := uc.service.GetConsumptionForSelectedDate(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

//
//func (uc SolarPanelController) GetValueFromLastMinute(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	CheckIfError(err, c)
//	graphData, err := uc.service.GetValueFromLastMinute(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
//		return
//	}
//	c.JSON(http.StatusOK, graphData)
//}

package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/controllers"
	"smarthome-back/services/devices/energetic"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type HomeBatteryController struct {
	service energetic.HomeBatteryService
}

func NewHomeBatteryController(db *sql.DB, influxDb influxdb2.Client, cacheService cache.CacheService) HomeBatteryController {
	return HomeBatteryController{service: energetic.NewHomeBatteryService(db, influxDb, cacheService)}
}

func (uc HomeBatteryController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (uc HomeBatteryController) GetConsumptionForLastHour(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	results := uc.service.GetConsumptionForLastHour(id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc HomeBatteryController) GetConsumptionForSelectedTime(c *gin.Context) {
	var input TimeInput
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetConsumptionForSelectedTime(input.Time, id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc HomeBatteryController) GetConsumptionForSelectedDate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
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

func (uc HomeBatteryController) GetStatusForSelectedTime(c *gin.Context) {
	var input TimeInput
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetStatusForSelectedTime(input.Time, id)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc HomeBatteryController) GetStatusForSelectedDate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
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
	results := uc.service.GetStatusForSelectedDate(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), id)
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

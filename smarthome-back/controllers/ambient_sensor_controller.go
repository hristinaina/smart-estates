package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/mqtt_client"
	"smarthome-back/services"
	"time"

	"github.com/gin-gonic/gin"
)

type AmbientSensorController struct {
	service services.AmbientSensorService
	mqtt    *mqtt_client.MQTTClient
}

func NewAmbientSensorController(db *sql.DB, mqtt *mqtt_client.MQTTClient) AmbientSensorController {
	return AmbientSensorController{service: services.NewAmbientSensorService(db), mqtt: mqtt}
}

type AmbientSensor struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func (uc AmbientSensorController) GetValueForHour(c *gin.Context) {
	results := mqtt_client.GetLastOneHourValues(uc.mqtt.GetInflux(), c.Param("id"))
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

	results := mqtt_client.GetValuesForSelectedTime(uc.mqtt.GetInflux(), input.Time, c.Param("id"))

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

	results := mqtt_client.GetValuesForDate(uc.mqtt.GetInflux(), startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), c.Param("id"))

	c.JSON(http.StatusOK, gin.H{"result": results})
}

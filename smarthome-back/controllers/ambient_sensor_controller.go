package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/mqtt_client"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
)

type AmbientSensorController struct {
	service services.AmbientSensorService
	mqtt    *mqtt_client.MQTTClient
}

func NewAmbientSensorController(db *sql.DB) AmbientSensorController {
	return AmbientSensorController{service: services.NewAmbientSensorService(db)}
}

type AmbientSensor struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

func (uc AmbientSensorController) GetValueForHour(c *gin.Context) {
	// id, err := strconv.Atoi(c.Param("id"))
	// CheckIfError(err, c)
	// device := uc.service.Get(id)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
	// 	return
	// }

	results := mqtt_client.GetLastOneHourValues(c.Param("id"))

	// fmt.Println(results)

	c.JSON(http.StatusOK, gin.H{"result": results})
}

package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/dtos/consumption_graph"
	"smarthome-back/services"
	"time"
)

type ElectricityController struct {
	service services.ElectricityService
}

func NewElectricityController(db *sql.DB, influxDb influxdb2.Client) ElectricityController {
	return ElectricityController{service: services.NewElectricityService(db, influxDb)}
}

// GetElectricityForSelectedTime used to get production OR consumption
func (uc ElectricityController) GetElectricityForSelectedTime(c *gin.Context) {
	var input consumption_graph.TimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetElectricityForSelectedTime(input.QueryType, input.Time, input.Type, input.SelectedOptions, "")
	c.JSON(http.StatusOK, gin.H{"result": results})
}

// GetElectricityForSelectedDate used to get production OR consumption
func (uc ElectricityController) GetElectricityForSelectedDate(c *gin.Context) {
	var input consumption_graph.DateInput

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
	results := uc.service.GetElectricityForSelectedDate(input.QueryType, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), input.Type, input.SelectedOptions, "")
	c.JSON(http.StatusOK, gin.H{"result": results})
}

// GetRatioForSelectedTime used to get ratio (production - consumption) and ed (electrical distribution)
func (uc ElectricityController) GetRatioForSelectedTime(c *gin.Context) {
	var input consumption_graph.TimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetRatioForSelectedTime(input.Time, input.Type, input.SelectedOptions, input.BatteryId)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

// GetRatioForSelectedDate used to get ratio (production - consumption) and ed (electrical distribution)
func (uc ElectricityController) GetRatioForSelectedDate(c *gin.Context) {
	var input consumption_graph.DateInput

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
	results := uc.service.GetRatioForSelectedDate(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), input.Type, input.SelectedOptions, input.BatteryId)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

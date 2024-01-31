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

type ConsumptionController struct {
	service services.ConsumptionService
}

func NewConsumptionController(db *sql.DB, influxDb influxdb2.Client) ConsumptionController {
	return ConsumptionController{service: services.NewConsumptionService(db, influxDb)}
}

func (uc ConsumptionController) GetConsumptionForSelectedTime(c *gin.Context) {
	var input consumption_graph.TimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetConsumptionForSelectedTime(input.QueryType, input.Time, input.Type, input.SelectedOptions)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc ConsumptionController) GetConsumptionForSelectedDate(c *gin.Context) {
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
	results := uc.service.GetConsumptionForSelectedDate(input.QueryType, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), input.Type, input.SelectedOptions)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc ConsumptionController) GetRatioForSelectedTime(c *gin.Context) {
	var input consumption_graph.TimeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	results := uc.service.GetRatioForSelectedTime(input.Time, input.Type, input.SelectedOptions)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

func (uc ConsumptionController) GetRatioForSelectedDate(c *gin.Context) {
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
	results := uc.service.GetRatioForSelectedDate(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), input.Type, input.SelectedOptions)
	c.JSON(http.StatusOK, gin.H{"result": results})
}

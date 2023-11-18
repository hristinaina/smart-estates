package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/services"
	"strconv"
)

type AirConditionerController struct {
	service services.AirConditionerService
}

func NewAirConditionerController(db *sql.DB) AirConditionerController {
	return AirConditionerController{service: services.NewAirConditionerService(db)}
}

func (uc AirConditionerController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

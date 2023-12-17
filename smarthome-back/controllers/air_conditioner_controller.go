package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/services"
	"strconv"

	"github.com/gin-gonic/gin"
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
	fmt.Println("KLIMA")
	fmt.Println(device)
	c.JSON(http.StatusOK, device)
}

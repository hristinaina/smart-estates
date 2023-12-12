package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/services"
	"strconv"
)

type SolarPanelController struct {
	service services.SolarPanelService
}

func NewSolarPanelController(db *sql.DB) SolarPanelController {
	return SolarPanelController{service: services.NewSolarPanelService(db)}
}

func (uc SolarPanelController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

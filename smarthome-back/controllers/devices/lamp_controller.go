package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/controllers"
	services "smarthome-back/services/devices"
	"strconv"
)

type LampController struct {
	service services.LampService
}

func NewLampController(db *sql.DB) LampController {
	return LampController{service: services.NewLampService(db)}
}

func (lc LampController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) GetAll(c *gin.Context) {
	lamps, err := lc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamps)
}

func (lc LampController) TurnOn(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.TurnOn(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) TurnOff(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.TurnOff(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

func (lc LampController) SetLightning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}
	level, err := strconv.Atoi(c.Param("level"))
	if controllers.CheckIfError(err, c) {
		return
	}
	lamp, err := lc.service.SetLightning(id, level)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lamp)
}

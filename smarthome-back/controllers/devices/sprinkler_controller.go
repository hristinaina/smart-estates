package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/controllers"
	services "smarthome-back/services/devices/outside"
	"strconv"
)

type SprinklerController struct {
	service services.SprinklerService
}

func NewSprinklerController(db *sql.DB, client influxdb2.Client) SprinklerController {
	return SprinklerController{service: services.NewSprinklerService(db, client)}
}

func (controller SprinklerController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return
	}

	sprinkler, err := controller.service.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(200, sprinkler)
}

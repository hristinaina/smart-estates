package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/controllers"
	"smarthome-back/services/devices/energetic"
	"strconv"
)

type EVChargerController struct {
	service energetic.EVChargerService
}

func NewEVChargerController(db *sql.DB, influxDb influxdb2.Client) EVChargerController {
	return EVChargerController{service: energetic.NewEVChargerService(db, influxDb)}
}

func (uc EVChargerController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	controllers.CheckIfError(err, c)
	device := uc.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

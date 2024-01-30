package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"smarthome-back/controllers"
	"smarthome-back/dtos"
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
	id := ExtractId(c)

	sprinkler, err := controller.service.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(200, sprinkler)
}

func (controller SprinklerController) GetAll(c *gin.Context) {
	sprinklers, err := controller.service.GetAll()
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, sprinklers)
}

func (controller SprinklerController) TurnOn(c *gin.Context) {
	id := ExtractId(c)

	sprinkler, err := controller.service.UpdateIsOn(id, true)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(204, sprinkler)
}

func (controller SprinklerController) TurnOff(c *gin.Context) {
	id := ExtractId(c)

	sprinkler, err := controller.service.UpdateIsOn(id, false)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(204, sprinkler)
}

func (controller SprinklerController) Delete(c *gin.Context) {
	id := ExtractId(c)

	isDeleted, err := controller.service.Delete(id)
	if (isDeleted == false) || (err != nil) {
		c.JSON(404, gin.H{"error": err})
		return
	}
	c.JSON(204, isDeleted)
}

func (controller SprinklerController) Add(c *gin.Context) {
	var dto dtos.DeviceDTO

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sprinkler, err := controller.service.Add(dto)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, sprinkler)
}

func (controller SprinklerController) AddSpecialMode(c *gin.Context) {
	id := ExtractId(c)

	var dto dtos.SprinklerSpecialModeDTO

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mode, err := controller.service.AddSpecialMode(id, dto)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, mode)
}

func (controller SprinklerController) GetSpecialModes(c *gin.Context) {
	id := ExtractId(c)

	modes, err := controller.service.GetSpecialModes(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, modes)
}

func (controller SprinklerController) DeleteSpecialMode(c *gin.Context) {
	id := ExtractId(c)

	isDeleted, err := controller.service.DeleteSpecialMode(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(204, isDeleted)
}

func (controller SprinklerController) GetSpecialMode(c *gin.Context) {
	id := ExtractId(c)

	mode, err := controller.service.GetSpecialMode(id)
	if controllers.CheckIfError(err, c) {
		return
	}

	c.JSON(200, mode)
}

func ExtractId(c *gin.Context) int {
	id, err := strconv.Atoi(c.Param("id"))
	if controllers.CheckIfError(err, c) {
		return -1
	}
	return id
}

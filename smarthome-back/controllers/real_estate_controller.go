package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/dtos"
	"smarthome-back/models"
	"smarthome-back/services"
	"strconv"
)

type RealEstateController struct {
	service services.RealEstateService
}

func NewRealEstateController(db *sql.DB) RealEstateController {
	return RealEstateController{service: services.NewRealEstateService(db)}
}

func (rec RealEstateController) GetAll(c *gin.Context) {
	realEstates, err := rec.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) GetAllByUserId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if CheckIfError(err, c) {
		return
	}
	realEstates, err := rec.service.GetByUserId(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if CheckIfError(err, c) {
		return
	}
	realEstate, err := rec.service.Get(id)
	if CheckIfError(err, c) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, realEstate)
}

func (rec RealEstateController) GetPending(c *gin.Context) {
	realEstates, err := rec.service.GetPending()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) ChangeState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if CheckIfError(err, c) {
		return
	}
	state, err := strconv.Atoi(c.Param("state"))
	if CheckIfError(err, c) {
		return
	}
	var reason dtos.DiscardRealEstate
	err = c.BindJSON(&reason)
	if CheckIfError(err, c) {
		return
	}
	realEstate, err := rec.service.ChangeState(id, state, reason.DiscardReason)
	if err != nil {
		fmt.Println("ERROR!!!")
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, realEstate)
}

func (rec RealEstateController) Add(c *gin.Context) {
	var estate models.RealEstate

	// convert json object to model real estate
	if err := c.BindJSON(&estate); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	estate, err := rec.service.Add(estate)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, estate)
}

func CheckIfError(err error, c *gin.Context) bool {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}
	return false
}

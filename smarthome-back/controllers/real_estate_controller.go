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
		fmt.Println("Error happened!")
		c.JSON(http.StatusBadRequest, "Error happened!")
	}
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) GetAllByUserId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	CheckIfError(err, c)
	realEstates := rec.service.GetAllByUserId(id)
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	realEstate, err := rec.service.Get(id)
	CheckIfError(err, c)
	c.JSON(http.StatusOK, realEstate)
}

func (rec RealEstateController) GetPending(c *gin.Context) {
	realEstates := rec.service.GetPending()
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) ChangeState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	state, err := strconv.Atoi(c.Param("state"))
	CheckIfError(err, c)
	var reason dtos.DiscardRealEstate
	err = c.BindJSON(&reason)
	CheckIfError(err, c)
	realEstate := rec.service.ChangeState(id, state, reason.DiscardReason)
	if realEstate.SquareFootage == 0 {
		c.JSON(http.StatusBadRequest, "Only pending real estates can be accepted/declined.")
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

	estate = rec.service.Add(estate)
	c.JSON(http.StatusOK, estate)
}

func CheckIfError(err error, c *gin.Context) bool {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, "Error happened!")
		return true
	}
	return false
}

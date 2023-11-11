package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smarthome-back/services"
	"strconv"
)

type RealEstateController struct {
	service services.RealEstateService
}

func NewRealEstateController(db *gorm.DB, database *sql.DB) RealEstateController {
	return RealEstateController{service: services.NewRealEstateService(db, database)}
}

func (rec RealEstateController) GetAll(c *gin.Context) {
	realEstates := rec.service.GetAll()
	if realEstates == nil {
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

func (rec RealEstateController) ChangeState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	CheckIfError(err, c)
	state, err := strconv.Atoi(c.Param("state"))
	CheckIfError(err, c)

	realEstate := rec.service.ChangeState(id, state)
	if realEstate.SquareFootage == 0 {
		c.JSON(http.StatusBadRequest, "Only pending real estates can be accepted/declined.")
		return
	}
	c.JSON(http.StatusOK, realEstate)

}

func CheckIfError(err error, c *gin.Context) bool {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, "Error happened!")
		return true
	}
	return false
}

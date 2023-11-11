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
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, "Error happened!")
	}
	realEstates := rec.service.GetAll(id)
	c.JSON(http.StatusOK, realEstates)
}

func (rec RealEstateController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, "Error happened!")
	}
	realEstate, err := rec.service.Get(id)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, "Error happened!")
	}
	c.JSON(http.StatusOK, realEstate)
}

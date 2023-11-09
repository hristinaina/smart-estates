package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smarthome-back/services"
)

type UserController struct {
	service services.UserService
}

func NewUserController(db *gorm.DB, database *sql.DB) UserController {
	return UserController{service: services.NewUserService(db, database)}
}

func (uc UserController) ListUsers(c *gin.Context) {
	users := uc.service.ListUsers()
	c.JSON(http.StatusOK, users)
}

func (uc UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.service.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc UserController) TestGetMethod(c *gin.Context) {
	uc.service.TestGetMethod()
}

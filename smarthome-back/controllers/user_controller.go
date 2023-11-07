package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"smarthome-back/services"
	"net/http"
)

type UserController struct {
	service services.UserService
}

func NewUserController(db *gorm.DB) UserController {
	return UserController{service: services.NewUserService(db)}
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

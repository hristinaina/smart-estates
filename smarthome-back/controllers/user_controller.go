package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/models"
	"smarthome-back/services"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service      services.UserService
	mail_service services.MailService
}

func NewUserController(db *sql.DB) UserController {
	return UserController{service: services.NewUserService(db), mail_service: services.NewMailService(db)}
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

var user = models.User{}

type ResetPassword struct {
	Email string `json:"email"`
}

func (uc UserController) SendResetPasswordEmail(c *gin.Context) {
	var input ResetPassword

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user.Email = input.Email

	// send mail
	expiration := time.Now().Add(time.Minute * 30)
	token, _ := uc.mail_service.GenerateToken(input.Email, expiration)
	uc.mail_service.SendVerifyEmail(input.Email, token)

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "A verification email has been sent"})
}

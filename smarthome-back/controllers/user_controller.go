package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"smarthome-back/services"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	mail_service services.MailService
	repo         repositories.UserRepository
}

func NewUserController(db *sql.DB, cacheService cache.CacheService) UserController {
	return UserController{mail_service: services.NewMailService(db), repo: repositories.NewUserRepository(db, &cacheService)}
}

var user = models.User{}

type VerifyEmail struct {
	Email string `json:"email"`
}

func (uc UserController) SendResetPasswordEmail(c *gin.Context) {
	var input VerifyEmail

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user.Email = input.Email

	// send mail
	expiration := time.Now().Add(time.Minute * 30)
	token, _ := uc.mail_service.GenerateToken(input.Email, expiration)

	// aync send mail
	go uc.mail_service.SendVerifyEmail(input.Email, token)

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "A verification email has been sent"})
}

type ResetPassword struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func (uc UserController) ResetPassword(c *gin.Context) {
	var input ResetPassword

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if !uc.mail_service.IsValidToken(input.Token) {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid token!"})
		return
	}

	// find user with user.Email
	user, err := uc.repo.GetUserByEmail(user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	// edit user
	err = uc.repo.ResetPassword(user.Email, string(hash))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save new password"})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "Successfully reset password!"})
}

package controllers

import (
	"database/sql"
	"net/http"
	"smarthome-back/models"
	"smarthome-back/repositories"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo repositories.UserRepository
}

func NewAuthController(db *sql.DB) AuthController {
	return AuthController{repo: repositories.NewUserRepository(db)}
}

// type LoginInput struct {
// 	Username string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

// func Login(c *gin.Context) {

// 	var input LoginInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	u := models.User{}

// 	u.Username = input.Username
// 	u.Password = input.Password

// 	token, err := models.LoginCheck(u.Username, u.Password)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})

// }

// request body
type RegisterInput struct {
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Picture  string `json:"picture" binding:"required"`
}

func (uc AuthController) Register(c *gin.Context) {

	// get values from req body
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	// create the user
	u := models.User{}

	u.Email = input.Username
	u.Password = string(hash)
	u.Name = input.Name
	u.Surname = input.Surname
	u.Picture = input.Picture

	if err := uc.repo.SaveUser(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

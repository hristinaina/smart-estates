package controllers

import (
	"database/sql"
	"net/http"
	"os"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo repositories.UserRepository
}

func NewAuthController(db *sql.DB) AuthController {
	return AuthController{repo: repositories.NewUserRepository(db)}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (uc AuthController) Login(c *gin.Context) {

	// get email and password from req body
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// looked up requested user by email
	user, err := uc.repo.GetUserByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// compare sent password with saved user hash password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	// send it back
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// request body
type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
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

	u.Email = input.Email
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

package controllers

import (
	"database/sql"
	"net/http"
	"os"
	"smarthome-back/cache"
	"smarthome-back/enumerations"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"smarthome-back/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo         repositories.UserRepository
	mail_service services.MailService
}

func NewAuthController(db *sql.DB, cacheService cache.CacheService) AuthController {
	return AuthController{repo: repositories.NewUserRepository(db, &cacheService), mail_service: services.NewMailService(db)}
}

// request body
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var u = models.User{}

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
		"sub":  user.Id,
		"role": strconv.Itoa(int(user.Role)),
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// request body
type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Role     int    `json:"role" binding:"required"`
}

func (uc AuthController) SendVerificationMail(c *gin.Context) {

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
	u.Email = input.Email
	u.Password = string(hash)
	u.Name = input.Name
	u.Surname = input.Surname
	u.Role = enumerations.IntToRole(input.Role)

	if _, err := uc.repo.GetUserByEmail(u.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account with that email already exists!"})
		return
	}

	// send mail
	expiration := time.Now().Add(time.Hour * 24)
	token, _ := uc.mail_service.GenerateToken(input.Email, expiration)

	// asyc send verification mail for reg
	go uc.mail_service.CreateVarificationMail(input.Email, input.Name, input.Surname, token)

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "check mail"})
}

func (uc AuthController) Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": user})
}

func (uc AuthController) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successful logout!"})
}

// request body
type ActivateAccount struct {
	Token string `json:"token" binding:"required"`
}

func (uc AuthController) ActivateAccount(c *gin.Context) {
	var input ActivateAccount

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if uc.mail_service.IsValidToken(input.Token) {
		// add to database
		uc.repo.SaveUser(u)
		c.JSON(http.StatusOK, gin.H{"message": "Valid token!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token!"})
	}
}

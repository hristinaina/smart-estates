package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"smarthome-back/enumerations"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"smarthome-back/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SuperAdminController interface {
	ResetPassword(c *gin.Context)
	AddAdmin(c *gin.Context)
	EditSuperAdmin(c *gin.Context)
}

type SuperAdminControllerImpl struct {
	db                 *sql.DB
	repo               repositories.UserRepository
	mail_service       services.MailService
	superadmin_service services.GenerateSuperadmin
}

func NewSuperAdminController(db *sql.DB) SuperAdminController {
	return &SuperAdminControllerImpl{db: db, repo: repositories.NewUserRepository(db), mail_service: services.NewMailService(db), superadmin_service: services.NewGenerateSuperAdmin(db)}
}

// password dto
type PasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func (sas *SuperAdminControllerImpl) ResetPassword(c *gin.Context) {
	// get values from req body
	var input PasswordInput

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

	// get superadmin the user
	admin, err := sas.repo.GetUserByEmail("admin")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin not found"})
		return
	}

	// edit superadmin
	if err := sas.repo.ResetSuperAdminPassword(string(hash), admin.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error from database!"})
		return
	}

	// remove admin.json file
	filePath := "admin.json"

	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Delete file error:", err)
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

type Admin struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

func (sas *SuperAdminControllerImpl) AddAdmin(c *gin.Context) {
	// receive admin mail from request
	var input Admin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// generate admin password
	password := sas.superadmin_service.GenerateRandomPassword(12)

	// save admin in database
	newAdmin := models.User{Email: input.Email, Password: sas.superadmin_service.HashPassword(password), Name: input.Name, Surname: input.Surname, Role: enumerations.ADMIN}
	err := sas.repo.SaveUser(newAdmin)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is already an account with the entered email"})
		return
	}

	// send admin mail with his password - async send mail
	go sas.mail_service.CreateAdminLoginRequest(input.Name, input.Surname, input.Email, password)

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "Successfully added admin"})
}

func (sas *SuperAdminControllerImpl) EditSuperAdmin(c *gin.Context) {
	var input Admin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	err := sas.repo.EditSuperAdmin(input.Name, input.Surname, input.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Database error"})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "Successfully edit profile"})
}

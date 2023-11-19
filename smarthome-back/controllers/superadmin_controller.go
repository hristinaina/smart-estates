package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"smarthome-back/repositories"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SuperAdminController interface {
	ResetPassword(c *gin.Context)
}

type SuperAdminControllerImpl struct {
	db   *sql.DB
	repo repositories.UserRepository
}

func NewSuperAdminController(db *sql.DB) SuperAdminController {
	return &SuperAdminControllerImpl{db: db, repo: repositories.NewUserRepository(db)}
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

	// todo remove admin.json file
	filePath := "admin.json"

	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Delete file error:", err)
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

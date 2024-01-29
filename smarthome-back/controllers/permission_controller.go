package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"smarthome-back/services"
	"time"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	userRepository       repositories.UserRepository
	permissionRepository repositories.PermissionRepository
	mailService          services.MailService
}

func NewPermissionController(db *sql.DB) PermissionController {
	return PermissionController{userRepository: repositories.NewUserRepository(db), permissionRepository: repositories.NewPermissionRepository(db), mailService: services.NewMailService(db)}
}

// request body
type GrantPermission struct {
	Emails         []string
	Devices        []int
	RealEstateId   int
	RealEstateName string
	User           string
}

var tokenEmail map[string]string

func (pc PermissionController) ReceiveGrantPermission(c *gin.Context) {
	var input GrantPermission
	currentUserFromCookie, _ := c.Get("user")

	currentUser := currentUserFromCookie.(*models.User)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	sendMail := false
	tokenEmail = make(map[string]string)

	for _, email := range input.Emails {
		if email != currentUser.Email {
			user, _ := pc.userRepository.GetUserByEmail(email)
			if user != nil {
				fmt.Println("POSTOJI KORISNIK")

				for _, deviceId := range input.Devices {
					isAlreadyExistPermission := pc.permissionRepository.IsPermissionAlreadyExist(input.RealEstateId, deviceId, email)
					if !isAlreadyExistPermission {
						sendMail = true
						pc.permissionRepository.AddInactivePermission(input.RealEstateId, deviceId, email)
					}
				}
				if sendMail {
					fmt.Println("SALJI MEJL")
					expiration := time.Now().Add(time.Minute * 30)
					token, _ := pc.mailService.GenerateToken(user.Email, expiration)
					go pc.mailService.PermissionMail(email, user.Name, input.User, input.RealEstateName, token)

					tokenEmail[token] = email
					sendMail = false
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Send mails!"})
}

func (pc PermissionController) VerifyAccount(c *gin.Context) {
	var input ActivateAccount

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if pc.mailService.IsValidToken(input.Token) {
		email, _ := tokenEmail[input.Token]
		fmt.Println(email)
		pc.permissionRepository.SetActivePermission(email)
		c.JSON(http.StatusOK, gin.H{"message": "Valid token!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token!"})
	}
}

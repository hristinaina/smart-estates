package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"smarthome-back/models"
	"smarthome-back/repositories"
	"smarthome-back/services"

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

func (pc PermissionController) ReceiveGrantPermission(c *gin.Context) {
	var input GrantPermission
	currentUserFromCookie, _ := c.Get("user")

	currentUser := currentUserFromCookie.(*models.User)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	sendMail := false

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
					// todo posalji mejl

					sendMail = false
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Send mails!"})
}

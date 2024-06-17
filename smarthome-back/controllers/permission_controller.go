package controllers

import (
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"net/http"
	"smarthome-back/cache"
	"smarthome-back/dtos"
	"smarthome-back/models"
	"smarthome-back/repositories"
	repositories2 "smarthome-back/repositories/devices"
	"smarthome-back/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	userRepository       repositories.UserRepository
	deviceRepository     repositories2.DeviceRepository
	permissionRepository repositories.PermissionRepository
	mailService          services.MailService
}

func NewPermissionController(db *sql.DB, client influxdb2.Client, cacheService cache.CacheService) PermissionController {
	return PermissionController{userRepository: repositories.NewUserRepository(db, &cacheService),
		permissionRepository: repositories.NewPermissionRepository(db), mailService: services.NewMailService(db),
		deviceRepository: repositories2.NewDeviceRepository(db, client, &cacheService)}
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

				for _, deviceId := range input.Devices {
					isAlreadyExistPermission := pc.permissionRepository.IsPermissionAlreadyExist(input.RealEstateId, deviceId, email)
					if !isAlreadyExistPermission {
						sendMail = true
						pc.permissionRepository.AddInactivePermission(input.RealEstateId, deviceId, email)
					}
				}
				if sendMail {
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

func (pc PermissionController) GetPermissionForRealEstate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	permissionsDTO := pc.permissionRepository.GetPermissionByRealEstate(id)
	var permissions []dtos.PermissionDTO

	for _, permission := range permissionsDTO {
		user, _ := pc.userRepository.GetUserByEmail(permission.UserEmail)
		permission.User = user.Name + " " + user.Surname
		permissions = append(permissions, permission)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func (pc PermissionController) DeletePermit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	var input []dtos.PermissionDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	for _, permission := range input {
		err = pc.permissionRepository.DeletePermission(id, permission)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete permissions"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted permissions"})
}

func (pc PermissionController) GetPermitRealEstate(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	user, e := pc.userRepository.GetUserById(userId)
	if e != nil {
		return
	}

	realEstates := pc.permissionRepository.GetPermitRealEstateByEmail(user.Email)

	c.JSON(http.StatusOK, realEstates)
}

func (pc PermissionController) GetDeviceForRealEstate(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return
	}

	user, e := pc.userRepository.GetUserById(userId)
	if e != nil {
		return
	}

	realEstateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	devices := pc.permissionRepository.GetDeviceForSharedRealEstate(user.Email, realEstateId)

	c.JSON(http.StatusOK, devices)
}

func (pc PermissionController) GetPermissionsForDevice(c *gin.Context) {
	deviceId, err := strconv.Atoi(c.Param("deviceId"))
	if err != nil {
		return
	}

	realEstateId, _ := pc.deviceRepository.GetRealEstateByDeviceId(deviceId)

	users := pc.permissionRepository.GetPermissionsForDevice(deviceId, realEstateId)

	c.JSON(http.StatusOK, users)
}

func (pc PermissionController) GetUsersForRealEstate(c *gin.Context) {
	deviceId, err := strconv.Atoi(c.Param("deviceId"))
	estateId, err := strconv.Atoi(c.Param("estateId"))
	if err != nil {
		return
	}

	users := pc.permissionRepository.GetAllUsersForRealEstate(deviceId, estateId)

	c.JSON(http.StatusOK, users)
}

package services

import (
	"gorm.io/gorm"
	"smarthome-back/models"
)

type UserService interface {
	ListUsers() []models.User
	GetUser(id string) (models.User, error)
}

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{db: db}
}

func (us *UserServiceImpl) ListUsers() []models.User {
	var users = []models.User{
		{Id: 1, Name: "Blue Train"},
		{Id: 2, Name: "Jeru"},
		{Id: 3, Name: "Sarah Vaughan and Clifford Brown"},
	}
	// var users []models.User
	// us.db.Find(&users)
	return users
}

func (us *UserServiceImpl) GetUser(id string) (models.User, error) {
	var user models.User
	result := us.db.First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}


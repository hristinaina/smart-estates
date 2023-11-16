package models

import "smarthome-back/enumerations"

type User struct {
	Id       int
	Email    string
	Password string
	Name     string
	Surname  string
	Picture  string
	Role     enumerations.Role
}

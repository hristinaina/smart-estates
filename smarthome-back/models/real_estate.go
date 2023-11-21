package models

import (
	"smarthome-back/enumerations"
)

type RealEstate struct {
	Id             int
	Name           string
	Type           enumerations.RealEstateType
	Address        string
	City           string // predefined list of the cities from database
	SquareFootage  float32
	NumberOfFloors int
	Picture        string // change this later (upload picture)
	State          enumerations.State
	User           int
	DiscardReason  string
}

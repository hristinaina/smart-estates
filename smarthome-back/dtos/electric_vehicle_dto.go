package dtos

type ElectricVehicleDTO struct {
	MaxCapacity     int
	CurrentCapacity float64
	StartCapacity   float64
	Active          bool
	PlugId          int
	Action          string
	Email           string
}

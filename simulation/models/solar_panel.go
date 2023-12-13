package models

type SolarPanel struct {
	Device         Device
	IsOn           bool    `json:"IsOn"`
	SurfaceArea    float64 `json:"SurfaceArea"`
	Efficiency     float64 `json:"Efficiency"`
	NumberOfPanels int     `json:"NumberOfPanels"`
	UserId         int     `json:"UserId"`
}

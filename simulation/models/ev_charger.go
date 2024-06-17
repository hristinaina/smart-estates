package models

type EVCharger struct {
	Device        Device
	ChargingPower float64 `json:"ChargingPower"`
	Connections   int     `json:"Connections"`
}

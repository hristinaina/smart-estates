package models

type AmbientSensor struct {
	Device           Device
	PowerSupply      PowerSupplyType
	PowerConsumption float64
}

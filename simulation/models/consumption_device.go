package models

import "simulation/enums"

type ConsumptionDevice struct {
	Device           Device
	PowerSupply      enums.PowerSupplyType
	PowerConsumption float64
}

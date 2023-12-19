package models

import "simulation/enums"


type ConsumptionDevice struct {
	PowerSupply      enums.PowerSupplyType
	PowerConsumption float64
}
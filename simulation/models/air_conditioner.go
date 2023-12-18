package models

type AirConditioner struct {
	Device         ConsumptionDevice
	MinTemperature float32
	MaxTemperature float32
	Mode           string
	SpecialMode    []SpecialMode
}

type SpecialMode struct {
	StartTime    string
	EndTime      string
	Mode         string
	Temperature  float32
	SelectedDays string
}

type ConsumptionDevice struct {
	Device           Device
	PowerSupply      PowerSupplyType
	PowerConsumption float64 // only if power supply is "home", otherwise it's null
}

type PowerSupplyType int

const (
	Autonomous PowerSupplyType = iota
	Home
)

type ReceiveValue struct {
	Mode      string
	Switch    bool
	Temp      float32
	Previous  string
	UserEmail string
}

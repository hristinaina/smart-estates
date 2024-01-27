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

// uvek ce biti auto, i previous ”
type SendValue struct {
	Mode   string
	Switch bool
	Temp   float32
}

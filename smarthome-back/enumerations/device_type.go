package enumerations

type DeviceType int

const (
	AmbientSensor DeviceType = iota
	AirConditioner
	WashingMachine
	Lamp
	VehicleGate
	Sprinkler
	SolarPanel
	BatteryStorage
	EVCharger
)

func (dt DeviceType) String() string {
	statuses := [...]string{"AmbientSensor", "AirConditioner", "WashingMachine", "Lamp", "VehicleGate", "Sprinkler",
		"SolarPanel", "BatteryStorage", "EVCharger"}

	// Check if the enum value is within the valid range
	if dt < 0 || int(dt) >= len(statuses) {
		return "Unknown"
	}

	return statuses[dt]
}

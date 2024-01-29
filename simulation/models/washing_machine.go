package models

type WashingMachine struct {
	Device        ConsumptionDevice
	Mode          []SupportMode
	ModeName      string
	ScheduledMode []ScheduledMode
}

type SupportMode struct {
	Id          int
	Name        string
	Duration    int
	Temperature string
}

type ScheduledMode struct {
	Id        int
	DeviceId  int
	StartTime string
	ModeId    int
}

type WMReceiveValue struct {
	Mode      string
	Switch    bool
	Previous  string
	UserEmail string
}

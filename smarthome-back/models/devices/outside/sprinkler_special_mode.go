package models

import "smarthome-back/enumerations"

type SprinklerSpecialMode struct {
	// TODO: vidjeti sa tasijom da njena SpecialModes naslijedi ovu klasu
	StartTime    string
	EndTime      string
	SelectedDays []enumerations.Days
}

func NewSprinklerSpecialMode(startTime, endTime string, selectedDays []enumerations.Days) SprinklerSpecialMode {
	return SprinklerSpecialMode{
		StartTime:    startTime,
		EndTime:      endTime,
		SelectedDays: selectedDays,
	}
}

func (s SprinklerSpecialMode) SelectedDaysToString() string {
	res := ""
	for _, day := range s.SelectedDays {
		res += day.String() + ","
	}

	return res
}

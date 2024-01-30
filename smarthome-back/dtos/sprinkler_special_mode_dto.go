package dtos

import (
	"smarthome-back/enumerations"
	models "smarthome-back/models/devices/outside"
)

type SprinklerSpecialModeDTO struct {
	StartTime    string
	EndTime      string
	SelectedDays string
}

func (mode SprinklerSpecialModeDTO) ToSprinklerSpecialMode() models.SprinklerSpecialMode {
	days, err := enumerations.ConvertStringsToEnumValues(mode.SelectedDays)
	if err != nil {
		return models.SprinklerSpecialMode{}
	}
	return models.NewSprinklerSpecialMode(mode.StartTime, mode.EndTime, days)
}

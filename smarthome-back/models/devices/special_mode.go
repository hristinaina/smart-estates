package models

import "time"

type SpecialMode struct {
	StartTime    time.Time
	EndTime      time.Time
	Mode         string
	Temperature  float32
	SelectedDays string
}

func NewSpecialMode(startTime, endTime time.Time, mode string, temperature float32, sc string) SpecialMode {
	return SpecialMode{
		StartTime:    startTime,
		EndTime:      endTime,
		Mode:         mode,
		Temperature:  temperature,
		SelectedDays: sc,
	}
}

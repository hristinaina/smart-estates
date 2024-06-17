package inside

type SpecialMode struct {
	StartTime    string
	EndTime      string
	Mode         string
	Temperature  float32
	SelectedDays string
}

func NewSpecialMode(startTime, endTime, mode string, temperature float32, sc string) SpecialMode {
	return SpecialMode{
		StartTime:    startTime,
		EndTime:      endTime,
		Mode:         mode,
		Temperature:  temperature,
		SelectedDays: sc,
	}
}

type SpecialModeDTO struct {
	Start        string
	End          string
	SelectedMode string
	Temperature  float32
	SelectedDays []string
}

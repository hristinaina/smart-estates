package dtos

type SpecialModeDTO struct {
	Start        string   `json:"start"`
	End          string   `json:"end"`
	SelectedMode string   `json:"selectedMode"`
	Temperature  float32  `json:"temperature"`
	SelectedDays []string `json:"selectedDays"`
}

package models

type HourlyData struct {
	Time                   []string  `json:"time"`
	DirectNormalIrradiance []float64 `json:"direct_normal_irradiance"`
}

type HourlyUnits struct {
	Time                   string `json:"time"`
	DirectNormalIrradiance string `json:"direct_normal_irradiance"`
}

type OpenMeteoResponse struct {
	Latitude             float64     `json:"latitude"`
	Longitude            float64     `json:"longitude"`
	GenerationTimeMS     float64     `json:"generationtime_ms"`
	UTCOffsetSeconds     int         `json:"utc_offset_seconds"`
	Timezone             string      `json:"timezone"`
	TimezoneAbbreviation string      `json:"timezone_abbreviation"`
	Elevation            float64     `json:"elevation"`
	HourlyUnits          HourlyUnits `json:"hourly_units"`
	Hourly               HourlyData  `json:"hourly"`
}

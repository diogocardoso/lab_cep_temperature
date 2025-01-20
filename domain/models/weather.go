package models

type Weather struct {
	CityName string  `json:"city_name"`
	TempC    float64 `json:"temp_C"`
	TempF    float64 `json:"temp_F"`
	TempK    float64 `json:"temp_K"`
}

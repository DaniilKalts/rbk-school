package dto

type WeatherResponse struct {
	City    string         `json:"city,omitempty"`
	Current CurrentWeather `json:"current"`
}
type CurrentWeather struct {
	Temperature2M       float64 `json:"temperature_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
}

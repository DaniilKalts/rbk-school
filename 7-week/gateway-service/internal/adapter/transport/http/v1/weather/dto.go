package weather

type Response struct {
	City                string  `json:"city"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Temperature         float64 `json:"temperature"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	WeatherCode         int     `json:"weather_code"`
}

package dto

type WeatherResponse struct {
	City           string     `json:"city,omitempty"`
	Conditions     Conditions `json:"conditions"`
	Recommendation string     `json:"recommendation"`
}

type Conditions struct {
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feelsLike"`
}

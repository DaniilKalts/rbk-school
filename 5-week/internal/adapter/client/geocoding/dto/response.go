package dto

type CoordsResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeocodingResults struct {
	Results []CoordsResponse `json:"results"`
}

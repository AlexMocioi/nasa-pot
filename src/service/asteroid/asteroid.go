package asteroid

import "strconv"

type ApiResponse struct {
	Links            map[string]string        `json:"links"`
	ElementCount     int                      `json:"element_count"`
	NearEarthObjects map[string][]interface{} `json:"near_earth_objects"`
}

type Specs struct {
	Name              string
	Id                string
	DiameterKmMin     float64
	DiameterKmMax     float64
	Hazardous         bool
	CloseApproachDate string
	MissDistanceKm    float64
}

func NewSpecs(jsonInput map[string]interface{}) Specs {
	asteroidData := Specs{}
	asteroidData.Id = jsonInput["id"].(string)
	asteroidData.Hazardous = jsonInput["is_potentially_hazardous_asteroid"].(bool)
	asteroidData.Name = jsonInput["name"].(string)
	estimatedDiameterKm := jsonInput["estimated_diameter"].(map[string]interface{})["kilometers"].(map[string]interface{})
	asteroidData.DiameterKmMin = estimatedDiameterKm["estimated_diameter_min"].(float64)
	asteroidData.DiameterKmMax = estimatedDiameterKm["estimated_diameter_max"].(float64)
	closeApproachData := jsonInput["close_approach_data"].([]interface{})[0].(map[string]interface{})
	asteroidData.CloseApproachDate, _ = closeApproachData["close_approach_date_full"].(string)
	asteroidData.MissDistanceKm, _ = strconv.ParseFloat(closeApproachData["miss_distance"].(map[string]interface{})["kilometers"].(string), 64)
	return asteroidData
}

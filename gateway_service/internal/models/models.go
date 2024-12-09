package models

type SearchRequest struct {
	Name          string `json:"name"`
	TypeOfFinding string `json:"type_of_finding"`
	Location      string `json:"location"`
}

type MakeAdsRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	TypeOfFinding string `json:"type_of_finding"`
	Location      string `json:"location"`
}

type ApplyRequest struct {
	Uuid string `json:"uuid"`
}

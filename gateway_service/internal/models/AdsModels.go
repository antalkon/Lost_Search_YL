package models

// WILL BE CHANGED NOT FINAL VARIANT

type Location struct {
	Country  string `json:"country"`
	City     string `json:"city"`
	District string `json:"district"`
}

type SearchRequest struct {
	Name          string   `json:"name"`
	TypeOfFinding string   `json:"type"`
	Location      Location `json:"location"`
}

type Finding struct {
	Name string
}

type SearchResponse struct {
	Finds []Finding `json:"finds"`
}

type SearchKafkaResponse struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Uuid        string   `json:"uuid"`
}

type MakeAdsRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TypeOfFinding string   `json:"type"`
	Location      Location `json:"location"`
}

type MakeAdsResponse struct {
	Uuid string `json:"uuid"`
}

type MakeAdsKafkaRequest struct {
	Login         string   `json:"login"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TypeOfFinding string   `json:"type"`
	Location      Location `json:"location"`
}

type MakeAdsKafkaResponse struct {
	Success bool `json:"success"`
}

type ApplyRequest struct {
	Uuid string `json:"uuid"`
}

type ApplyResponse struct {
	Login string `json:"user_login"`
}

type ApplyKafkaRequest struct {
	Uuid string `json:"find_uuid"`
}

type ApplyKafkaResponse struct {
	Uuid string `json:"uuid"`
}

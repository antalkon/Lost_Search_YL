package models

// WILL BE CHANGED NOT FINAL VARIANT
type KafkaRequest struct {
	RequestId string `json:"request_id"`
	Service   string `json:"service"`
	Action    string `json:"action"`
	Data      string `json:"data"`
}

type KafkaResponse struct {
	RequestId string `json:"request_id"`
	Service   string `json:"service"`
	Status    string `json:"action"`
	Data      string `json:"data"`
}

package models

// WILL BE CHANGED NOT FINAL VARIANT
type NotifyRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

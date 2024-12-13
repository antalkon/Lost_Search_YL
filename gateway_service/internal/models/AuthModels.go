package models

// WILL BE CHANGED NOT FINAL VARIANT
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"data"`
}

type RegisterKafkaResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginKafkaResponse struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	Success bool `json:"success"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid bool `json:"valid"`
}

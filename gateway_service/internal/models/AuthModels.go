package models

// WILL BE CHANGED NOT FINAL VARIANT
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"data"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

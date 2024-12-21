package repository

import (
	"encoding/json"
)

type Repository interface {
	AddFind(AddReq) (AddResp, error)
	GetFind(GetReq) (GetResp, error)
	RespondToFind(RespondReq) (RespondResp, error)
}

type location struct {
	Country  string `json:"country" db:"country"`
	City     string `json:"city" db:"city"`
	District string `json:"district" db:"district"`
}

type find struct {
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Type        string   `json:"type" db:"type"`
	Location    location `json:"location"`
}

type AddReq struct {
	find
	UserLogin string `json:"user_login" db:"user_login"`
}

type AddResp struct {
	FindUUID string `json:"find_uuid"`
}

type GetReq struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location location `json:"location"`
}

type GetResp struct {
	Finds []AddReq `json:"finds"`
}

type RespondReq struct {
	FindUUID string `json:"find_uuid"`
}

type RespondResp struct {
	UserLogin string `json:"user_login"`
}

func (r *GetResp) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GetResp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func AddReqMok() AddReq {
	return AddReq{
		UserLogin: "egorc",
		find: find{
			Name:        "телефон",
			Description: "потерялся",
			Type:        "техника",
			Location: location{
				Country:  "Россия",
				City:     "Иваново",
				District: "Ленинский",
			},
		},
	}
}

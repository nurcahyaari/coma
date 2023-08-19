package dto

import "time"

type ResponseValidateKey struct {
	Valid bool `json:"valid"`
}

type ResponseExtractedToken struct {
	UserId    string    `json:"userId"`
	ExpiredAt time.Time `json:"expiredAt"`
	UserType  string    `json:"userType"`
}

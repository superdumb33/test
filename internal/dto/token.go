package dto

import "github.com/google/uuid"

type AuthorizeUserRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type AuthorizeUserResponse struct {
	AccessToken string`json:"access_token"`
	RefreshToken string`json:"refresh_token"`
}
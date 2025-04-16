package dto

import "github.com/google/uuid"

type AuthorizeUserRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type AuthorizeUserResponse struct {
	AccessToken string`json:"access_token"`
	RefreshToken string`json:"refresh_token"`
}

type RefreshTokensRequest struct{
	UUID uuid.UUID `json:"uuid"`
	AccessToken string`json:"access_token"`
	RefreshToken string`json:"refresh_token"`
}
type RefreshTokensResponse struct{
	AccessToken string`json:"access_token"`
	RefreshToken string`json:"refresh_token"`
}
package services

import (
	"encoding/base64"
	"rest-service/internal/auth"
	"rest-service/internal/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(*entities.User) error
}

type UserService struct {
	repo UserRepository
}
//nice naming =)
type Tokens struct {
	AccessToken string
	RefreshToken string
}

func NewUserService (repo UserRepository) *UserService{
	return &UserService{repo: repo}
}

func (us *UserService) Authorize (userID uuid.UUID, userIP string) (Tokens, error) {
	tokenID := uuid.New()	

	accesToken, err := auth.GenerateAccessToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err
	}
	refreshToken, err := auth.GenerateRefreshToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err
	}

	hash, err := auth.GenerateBCryptHash(refreshToken)
	if err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err
	}

	if err := us.repo.CreateUser(&entities.User{ID:userID, RefreshToken: string(hash)}); err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err 
	}
	base64EncodedRefreshToken := base64.StdEncoding.EncodeToString(hash)

	return Tokens{
		AccessToken: accesToken,
		RefreshToken: base64EncodedRefreshToken,
	}, nil
}
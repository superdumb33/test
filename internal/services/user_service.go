package services

import (
	"rest-service/internal/auth"
	"rest-service/internal/entities"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err
	}

	if err := us.repo.CreateUser(&entities.User{ID: userID, RefreshToken: string(refreshTokenHash)}); err != nil {
		return Tokens{
			AccessToken: "",
			RefreshToken: "",
		}, err
	}

	return Tokens{
		AccessToken: accesToken,
		RefreshToken: refreshToken,
	}, nil

}
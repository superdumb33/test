package services

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"rest-service/internal/auth"
	"rest-service/internal/entities"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(context.Context, *entities.User) error
	GetUser(context.Context, string) (*entities.User, error)
	UpdateUser(context.Context, *entities.User) error
}

type UserService struct {
	repo UserRepository
}

// nice naming =)
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Authorize(ctx context.Context, userID uuid.UUID, userIP string) (Tokens, error) {
	tokenID := uuid.New()

	accesToken, err := auth.GenerateAccessToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}
	refreshToken, err := auth.GenerateRefreshToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}

	hash, err := auth.GenerateBCryptHash(refreshToken)
	if err != nil {
		return Tokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}

	if err := us.repo.CreateUser(ctx, &entities.User{ID: userID, RefreshToken: string(hash)}); err != nil {
		return Tokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}
	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	return Tokens{
		AccessToken:  accesToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}


func (us *UserService) Refresh(ctx context.Context, accessToken, refreshToken string, userIP string) (Tokens, error) {
	token, err := auth.ParseJWTToken(accessToken)
	if err != nil || !token.Valid {
		return Tokens{}, errors.New("unauthorized")
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	tokenID := claims["token_id"].(string)

	user, err := us.repo.GetUser(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}

	rawRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	refreshTokenString := string(rawRefreshToken)
	if err := auth.VerifyRefreshToken(refreshTokenString, user.RefreshToken); err != nil {
		return Tokens{}, err
	}

	parsedRefreshtoken, err := auth.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return Tokens{}, err
	}
	if parsedRefreshtoken.TokenID != tokenID {
		return Tokens{}, errors.New("mismatched token ids")
	}
	if parsedRefreshtoken.UserIP != userIP {
		log.Println("mismatched ip adress")
	}

	newTokenID := uuid.New()

	newAccessToken, err := auth.GenerateAccessToken(user.ID.String(), userIP, newTokenID.String())
	if err != nil {
		return Tokens{}, err
	}
	newRefreshToken, err := auth.GenerateRefreshToken(user.ID.String(), userIP, newTokenID.String())
	if err != nil {
		return Tokens{}, err
	}
	hash, err := auth.GenerateBCryptHash(newRefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if err := us.repo.UpdateUser(ctx, &entities.User{ID: user.ID, RefreshToken: string(hash)}); err != nil {
		return Tokens{}, err
	}

	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))


	return Tokens{
		AccessToken: newAccessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

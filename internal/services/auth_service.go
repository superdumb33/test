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

type AuthRepository interface {
	CreateUser(context.Context, *entities.User) error
	GetUser(ctx context.Context, userID string) (*entities.User, error)
	UpdateUser(context.Context, *entities.User) error
}

type SMTPClient interface {
	Send(to []string, subject, body string) error
}

type AuthService struct {
	repo       AuthRepository
	smtpclient SMTPClient
}

// nice naming =)
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

var (
	GenerateAccessToken  = auth.GenerateAccessToken
	GenerateRefreshToken = auth.GenerateRefreshToken
	GenerateBCryptHash   = auth.GenerateBCryptHash
	VerifyRefreshToken   = auth.VerifyRefreshToken
	ParseJWTToken        = auth.ParseJWTToken
	ParseRefreshToken    = auth.ParseRefreshToken
)

func NewUserService(repo AuthRepository, smtpClient SMTPClient) *AuthService {
	return &AuthService{repo: repo, smtpclient: smtpClient}
}

func (as *AuthService) Authorize(ctx context.Context, userID uuid.UUID, userIP string) (Tokens, error) {
	tokenID := uuid.New()

	accessToken, err := GenerateAccessToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{}, err
	}
	refreshToken, err := GenerateRefreshToken(userID.String(), userIP, tokenID.String())
	if err != nil {
		return Tokens{}, err
	}

	hash, err := GenerateBCryptHash(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if err := as.repo.CreateUser(ctx, &entities.User{ID: userID, RefreshToken: string(hash)}); err != nil {
		return Tokens{}, err
	}
	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

func (as *AuthService) Refresh(ctx context.Context, accessToken, refreshToken string, userIP string) (Tokens, error) {
	token, err := ParseJWTToken(accessToken)
	if err != nil || !token.Valid {
		return Tokens{}, errors.New("unauthorized")
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	tokenID := claims["token_id"].(string)

	user, err := as.repo.GetUser(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}

	rawRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	refreshTokenString := string(rawRefreshToken)
	if err := VerifyRefreshToken(refreshTokenString, user.RefreshToken); err != nil {
		return Tokens{}, err
	}

	parsedRefreshtoken, err := ParseRefreshToken(refreshTokenString)
	if err != nil {
		return Tokens{}, err
	}
	if parsedRefreshtoken.TokenID != tokenID {
		return Tokens{}, errors.New("mismatched token ids")
	}
	if parsedRefreshtoken.UserIP != userIP {
		if err := as.smtpclient.Send([]string{"exampleuser@mail.com"}, "IP Address mismatch warning", "Detected an attempt to authorize from another location"); err != nil {
			log.Println(err)
		}
	}

	newTokenID := uuid.New()

	newAccessToken, err := GenerateAccessToken(user.ID.String(), userIP, newTokenID.String())
	if err != nil {
		return Tokens{}, err
	}
	newRefreshToken, err := GenerateRefreshToken(user.ID.String(), userIP, newTokenID.String())
	if err != nil {
		return Tokens{}, err
	}
	hash, err := GenerateBCryptHash(newRefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if err := as.repo.UpdateUser(ctx, &entities.User{ID: user.ID, RefreshToken: string(hash)}); err != nil {
		return Tokens{}, err
	}

	encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))

	return Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: encodedRefreshToken,
	}, nil
}

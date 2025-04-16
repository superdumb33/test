package services_test

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"rest-service/internal/entities"
	"rest-service/internal/services"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type MockAuthRepo struct {
	Users map[string]string
}

func (mr *MockAuthRepo) CreateUser(ctx context.Context, user *entities.User) error {
	mr.Users[user.ID.String()] = user.RefreshToken
	return nil
}

func (mr *MockAuthRepo) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	token := mr.Users[userID]
	id, _ := uuid.Parse(userID)
	return &entities.User{ID: id, RefreshToken: token}, nil
}

func (mr *MockAuthRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	return nil
}

func TestAuthorize(t *testing.T) {
	authService := services.NewUserService(&MockAuthRepo{Users: make(map[string]string)})
	userID := uuid.New()
	userIP := "127.0.0.1"

	t.Run("Valid data", func(t *testing.T) {
		tokens, err := authService.Authorize(context.Background(), userID, userIP)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if tokens.AccessToken == "" || tokens.RefreshToken == "" {
			t.Error("epmty token being returned")
		}
	})

	t.Run("Error while generating access token", func(t *testing.T) {
		originalGenerateAccesTokenFunc := services.GenerateAccessToken
		defer func() {
			services.GenerateAccessToken = originalGenerateAccesTokenFunc
		}()

		services.GenerateAccessToken = func(userID, userIP, tokenID string) (string, error) {
			return "", errors.New("test generating access error")
		}

		expectedError := "test generating access error"
		_, err := authService.Authorize(context.Background(), userID, userIP)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("Error while generating refresh token", func(t *testing.T) {
		originalGenerateRefreshTokenFunc := services.GenerateRefreshToken
		defer func() {
			services.GenerateRefreshToken = originalGenerateRefreshTokenFunc
		}()

		services.GenerateRefreshToken = func(userID, userIP, tokenID string) (string, error) {
			return "", errors.New("test generating refresh token error")
		}

		expectedError := "test generating refresh token error"
		_, err := authService.Authorize(context.Background(), userID, userIP)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("Error while generating bcrypt hash", func(t *testing.T) {
		originalGenerateBcryptHashFunc := services.GenerateBCryptHash
		defer func() {
			services.GenerateBCryptHash = originalGenerateBcryptHashFunc
		}()

		services.GenerateBCryptHash = func(token string) ([]byte, error) {
			return nil, errors.New("test hashing error")
		}

		expectedError := "test hashing error"
		_, err := authService.Authorize(context.Background(), userID, userIP)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})
}

func TestRefresh(t *testing.T) {
	authService := services.NewUserService(&MockAuthRepo{Users: make(map[string]string)})
	userID := uuid.New()

	tokens, err := authService.Authorize(context.Background(), userID, "127.0.0.1")
	if err != nil {
		t.Errorf("got an error during authorizing; err: %v", err)
	}

	t.Run("Valid tokens", func(t *testing.T) {
		refreshedTokens, err := authService.Refresh(context.Background(), tokens.AccessToken, tokens.RefreshToken, "127.0.0.1")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if refreshedTokens.AccessToken == "" || refreshedTokens.RefreshToken == "" {
			t.Errorf("empty tokens")
		}
	})

	t.Run("Invalid access token", func(t *testing.T) {
		_, err := authService.Refresh(context.Background(), "abc", tokens.RefreshToken, "127.0.0.1")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected error %v, got %v", "unauthorized", err)
		}
	})

	t.Run("Invalid refresh token", func(t *testing.T) {
		_, err := authService.Refresh(context.Background(), tokens.AccessToken, "abc", "127.0.0.1")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Error while verifying refresh token", func(t *testing.T) {
		originalVerifyRefreshTokenFunc := services.VerifyRefreshToken
		defer func() {
			services.VerifyRefreshToken = originalVerifyRefreshTokenFunc
		}()

		services.VerifyRefreshToken = func(refreshToken, storedHash string) error {
			return errors.New("test verifying refresh token error")
		}

		expectedError := "test verifying refresh token error"
		_, err := authService.Refresh(context.Background(), tokens.AccessToken, tokens.RefreshToken, "127.0.0.1")
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("Mismatched token id", func(t *testing.T) {
		originalVerifyRefreshTokenFunc := services.VerifyRefreshToken
		defer func() {
			services.VerifyRefreshToken = originalVerifyRefreshTokenFunc
		}()

		services.VerifyRefreshToken = func (refreshToken, storedHash string) error {
			return nil
		}

		rawToken, _ := base64.StdEncoding.DecodeString(tokens.RefreshToken)
		parts := strings.Split(string(rawToken), ":")
		parts[3] = "invalidTokenID"
		refreshToken := strings.Join(parts, ":")

		encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
		_, err = authService.Refresh(context.Background(), tokens.AccessToken, encodedRefreshToken, "127.0.0.1")
		if err == nil || err.Error() != "mismatched token ids" {
			log.Println(err)
			t.Errorf("expected error, got nil")
		}
	})
}

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAccessToken(userID, userIP, tokenID string) (string, error) {
	claims := jwt.MapClaims{
		"id": userID,
		"ip": userIP,
		"exp": time.Now().Add(time.Hour),
		"token_id": tokenID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken (userID, userIP, tokenID string) (string, error) {
	var randBytes = make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	randString := hex.EncodeToString(randBytes)

	token := fmt.Sprint(randString + ":" + userID + ":" + userIP + ":" + tokenID)

	return token, nil
}

//accepts raw token, calculates it's sha256 hash, to ensure it meets bcrypt's maximum lenght of 72 bytes, then calculates bcrypt hash
func GenerateBCryptHash (token string) ([]byte, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	sha256Hash := sha256.Sum256([]byte(token))

	return bcrypt.GenerateFromPassword(sha256Hash[:], bcrypt.DefaultCost)
}
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	RandomBytes string
	UserID string
	UserIP string
	TokenID string
}

func GenerateAccessToken(userID, userIP, tokenID string) (string, error) {
	claims := jwt.MapClaims{
		"id": userID,
		"ip": userIP,
		"exp": time.Now().Add(time.Hour).Unix(),
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

//accepts raw token, calculates it's sha256 hash, to ensure it meets bcrypt's maximum lenght of 72 bytes, returns bcrypt hash 
func GenerateBCryptHash (token string) ([]byte, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	sha256Hash := sha256.Sum256([]byte(token))

	return bcrypt.GenerateFromPassword(sha256Hash[:], bcrypt.DefaultCost)
}

func ParseJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, errors.New("unprocessable signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func ParseRefreshToken (tokenString string) (RefreshToken, error) {
	if tokenString == "" {
		return RefreshToken{}, errors.New("empty token string")
	}

	data := strings.Split(tokenString, ":")

	return RefreshToken{
		RandomBytes: data[0],
		UserID: data[1],
		UserIP: data[2],
		TokenID: data[3],
	}, nil
}

//accepts refreshToken, calculates it's sha256 hash, then comparing it's bcrypt hash with storedHash; returns nil if hash matches
func VerifyRefreshToken (refreshToken, storedHash string) error {
	shaHash := sha256.Sum256([]byte(refreshToken))

	return bcrypt.CompareHashAndPassword([]byte(storedHash), shaHash[:])
}
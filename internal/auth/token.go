package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	rand.Read(randBytes)
	randString := hex.EncodeToString(randBytes)

	token := fmt.Sprint(randString + ":" + userID + ":" + userIP + ":" + tokenID)

	return base64.RawStdEncoding.EncodeToString([]byte(token)), nil
}
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}
	uuid, _ := uuid.Parse(token.Claims.(*jwt.RegisteredClaims).Subject)
	return uuid, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	tokenstring := headers.Get("Authorization")

	tokenstring, ok := strings.CutPrefix(tokenstring, "Bearer ")
	if !ok || len(tokenstring) == 0 {
		return "", errors.New("error getting jwt token")
	}
	return tokenstring, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	rand.Read(token)

	return hex.EncodeToString(token), nil

}

func GetAPIKey(headers http.Header) (string, error) {

	apiKey := headers.Get("Authorization")
	apiKey, ok := strings.CutPrefix(apiKey, "ApiKey ")
	if !ok || len(apiKey) == 0 {
		return "", errors.New("error getting apikey")
	}
	return apiKey, nil
}

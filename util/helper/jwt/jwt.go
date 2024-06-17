package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signedKey = []byte("test")

func CreateRefreshToken(email string, tokenExpiry time.Duration) (string, *Payload, error) {
	return createToken(email, tokenExpiry)
}

func CreateAccessToken(email string, tokenExpiry time.Duration) (string, *Payload, error) {
	return createToken(email, tokenExpiry)
}

func createToken(email string, tokenExpiry time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(email, tokenExpiry)
	if err != nil {
		return "", nil, err // Added signed with error handling
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Create token with signed
	tokenString, err := token.SignedString(signedKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, payload, nil
}

func VerifyToken(tokenString string) (*Payload, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signedKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	emailInterface, ok := claims["Email"]
	if !ok {
		return nil, fmt.Errorf("email claim not found in token")
	}

	email, ok := emailInterface.(string)
	if !ok {
		return nil, fmt.Errorf("email claim is not a string")
	}

	payload := &Payload{
		Email: email,
	}

	return payload, nil
}
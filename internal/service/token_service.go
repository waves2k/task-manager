package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ParseToken(token string) (interface{}, error)
}

type jwtService struct {
	SecretKey string
}

func NewJWTService(key string) *jwtService {
	return &jwtService{
		SecretKey: key,
	}
}

func (s *jwtService) GenerateToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *jwtService) ParseToken(token string) (interface{}, error) {
	unparsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if err != nil || !unparsedToken.Valid {
		return nil, fmt.Errorf("Invalid authorization header format")
	}
	userId, ok := unparsedToken.Claims.(jwt.MapClaims)["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid authorization header format")
	}
	if exp, ok := unparsedToken.Claims.(jwt.MapClaims)["exp"].(int64); ok {
		expirationTime := time.Unix(exp, 0)
		if time.Now().After(expirationTime) {
			return nil, err
		}
	}
	return userId, nil
}

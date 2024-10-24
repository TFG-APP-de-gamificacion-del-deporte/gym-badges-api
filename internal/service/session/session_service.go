package session_service

import (
	configs "gym-badges-api/config/gym-badges-server"
	"gym-badges-api/internal/constants"
	customErrors "gym-badges-api/internal/custom-errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewSessionService() ISessionService {
	return &sessionService{}
}

type sessionService struct {
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (sessionService) GenerateSession(username string) (string, error) {

	expirationTime := time.Now().Add(time.Duration(configs.Basic.SessionDuration) * time.Second)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.Basic.JWTKey))
	if err != nil {
		return constants.EmptyString, err
	}

	return tokenString, nil
}

func (sessionService) GetUserFromToken(token string) (string, error) {
	var claims Claims
	tk, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Basic.JWTKey), nil
	})

	if err != nil || !tk.Valid {
		return "", customErrors.BuildUnauthorizedError("Invalid token")
	}

	return claims.Username, nil
}

func (s sessionService) ValidateSession(username string, sessionID string) error {
	tokenUsername, err := s.GetUserFromToken(sessionID)

	if err != nil {
		return err
	}

	if username != tokenUsername {
		return customErrors.BuildUnauthorizedError("Invalid token")
	}

	return nil
}

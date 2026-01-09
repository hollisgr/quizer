package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey string
}

func NewManager(secret string) *Manager {
	return &Manager{secretKey: secret}
}

func (m *Manager) CreateToken(userID int, login string) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userID,
		"login":   login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}
	return int(userID), nil
}

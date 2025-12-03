package jwt

import (
	"context"
	"log"
	"quizer_server/internal/config"
	"quizer_server/internal/model"
	"quizer_server/internal/service/user"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	CreateToken(ctx context.Context, req model.JwtRequest) model.JwtResponce
	ParseToken(token string, key string) (*jwt.Token, error)
	IDFromToken(tokenStr string) int
}

type jwtService struct {
	service user.Service
	cfg     *config.Config
}

// New creates a new instance of JwtService with the provided user repository.
func New(s user.Service) Service {
	return &jwtService{
		service: s,
		cfg:     config.GetConfig(),
	}
}

// IDFromToken extracts the user ID from a JWT token.
// It parses the token, retrieves the ID claim, and returns the parsed UUID.
func (js *jwtService) IDFromToken(tokenStr string) int {
	resp := 0
	token, err := js.ParseToken(tokenStr, js.cfg.Jwt.SecretKey)

	if err != nil {
		log.Println("jwt_service: parse token err: ", err)
		return resp
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		log.Println("jwt_service: get user_id claims err")
		return resp
	}

	user_id, ok := claims["user_id"].(float64)

	if !ok {
		log.Println("jwt_service: get user_id err")
		return resp
	}
	return int(user_id)
}

// CreateTokenPair generates a pair of access and refresh tokens for a user.
// It retrieves the user record by login, creates the access token with limited duration and refresh token with extended validity period,
// and returns them in a structured response.
func (js *jwtService) CreateToken(ctx context.Context, req model.JwtRequest) model.JwtResponce {
	user, err := js.service.UserByLogin(ctx, req.Login)

	if err != nil {
		log.Println("jwt_service get user err: ", err)
		return model.JwtResponce{}
	}

	aToken := js.createAccessToken(user.Id, user.Login)

	return model.JwtResponce{
		AccessToken: aToken,
		UserID:      user.Id,
	}
}

// ParseToken verifies and parses a JWT token using the provided secret key.
func (js *jwtService) ParseToken(tokenString string, key string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
}

func (js *jwtService) createAccessToken(userId int, login string) string {

	payload := jwt.MapClaims{
		"user_id": userId,
		"login":   login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, _ := token.SignedString([]byte(js.cfg.Jwt.SecretKey))
	return t
}

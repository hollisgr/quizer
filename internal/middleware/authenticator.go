package middleware

import (
	"fmt"
	"net/http"
	"quizer_server/internal/config"
	"quizer_server/internal/service/jwt"
	"quizer_server/internal/service/user"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserAuthenticator interface {
	Authorization() gin.HandlerFunc
}

type userAuthenticator struct {
	userService user.Service
	jwtService  jwt.Service
}

// NewUserAuthenticator creates a new instance of UserAuthenticator with provided dependencies.
// It takes a Repository interface for user-related operations and JwtService for JWT token management.
func NewUserAuthenticator(u user.Service, js jwt.Service) UserAuthenticator {
	return &userAuthenticator{
		userService: u,
		jwtService:  js,
	}
}

// Authorization implements a middleware handler for authentication purposes.
// It extracts a token from the request header, validates its authenticity against the JWT secret key,
// and stores the valid token in the context before proceeding to next handler.
func (a *userAuthenticator) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := parseTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			sendError(c, http.StatusUnauthorized, fmt.Sprint(err))
			return
		}

		_, err = a.jwtService.ParseToken(token, config.GetConfig().Jwt.SecretKey)

		if err != nil {
			sendError(c, http.StatusUnauthorized, "Access denied, invalid access token")
			return
		}

		c.Set("access_token", token)
		c.Next()
	}
}

// parseTokenFromHeader extracts the actual token value from the 'Authorization' header.
// It expects the format 'Bearer <token>' and returns either the extracted token or an appropriate error message.
func parseTokenFromHeader(header string) (string, error) {
	if header == "" {
		return header, fmt.Errorf("access denied, Authorization header required")
	}
	headerArr := strings.Split(header, " ")
	if headerArr[0] != "Bearer" {
		return header, fmt.Errorf("access denied, Bearer authorization required")
	}
	return headerArr[1], nil
}

// sendError sends an error response back to the client with a specific HTTP status code and custom error message.
func sendError(c *gin.Context, code int, message any) {
	c.AbortWithStatusJSON(code, gin.H{
		"success": false,
		"message": fmt.Sprint(message),
	})
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/api"
	"github.com/waves2k/task-manager/internal/service"
)

const (
	authorizationHeader = "Authorization"
)

type AuthMiddleware struct {
	tokenService service.TokenService
}

func NewAuthMiddleware(tokenService service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			api.NewErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		if tokenString == "" || tokenString == authHeader {
			api.NewErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}
		userId, err := m.tokenService.ParseToken(tokenString)
		if err != nil {
			api.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", userId)
		c.Next()
	}
}

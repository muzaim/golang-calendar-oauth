package middleware

import (
	"net/http"
	"strings"

	"golang-test/pkg/auth"
	"golang-test/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey    = "userID"
	CtxUserEmailKey = "userEmail"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired access token", nil)
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUserEmailKey, claims.Email)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(CtxUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}

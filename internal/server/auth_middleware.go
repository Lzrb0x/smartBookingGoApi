package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authenticatedUserIDKey = "authenticated_user_id"

type accessTokenClaims struct {
	TokenType string `json:"typ"`
	UserID    int64  `json:"uid"`
	jwt.RegisteredClaims
}

func authMiddleware(jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de acesso obrigatório"})
			return
		}

		tokenValue, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(tokenValue) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de acesso inválido"})
			return
		}

		claims := &accessTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid || claims.TokenType != "access" || claims.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de acesso inválido"})
			return
		}

		c.Set(authenticatedUserIDKey, claims.UserID)
		c.Next()
	}
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints/response"
	"github.com/modhanami/boinger/services"
	"net/http"
	"strings"
)

const (
	UserClaimsKey = "userClaims"
)

func MakeVerifyJWTUserTokenMiddleware(s services.UserTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(header, " ")
		tokenType, token := tokenParts[0], tokenParts[1]
		if tokenType != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := s.Verify(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponseFromError(err))
			c.Abort()
			return
		}

		c.Set(UserClaimsKey, claims)
	}
}

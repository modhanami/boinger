package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/services"
	"net/http"
)

func MakeVerifyJWTUserTokenMiddleware(s services.UserTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenCookie, err := c.Request.Cookie(endpoints.AuthTokenCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, endpoints.ErrorResponseFromError(err))
			c.Abort()
			return
		}

		claims, err := s.Verify(authTokenCookie.Value)
		if err != nil {
			c.JSON(http.StatusUnauthorized, endpoints.ErrorResponseFromError(err))
			c.Abort()
			return
		}

		c.Set(endpoints.UserIdKey, claims.ID)
		c.Set(endpoints.UserClaimsKey, claims)
	}
}

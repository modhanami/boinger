package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/services"
	"net/http"
	"strings"
)

var UserIdKey = "userId"

func MakeVerifyJWTUserTokenMiddleware(s services.UserTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawAuthHeader := c.GetHeader("Authorization")
		if rawAuthHeader == "" {
			c.JSON(http.StatusUnauthorized, endpoints.NewErrorResponse("No Authorization header"))
			c.Abort()
			return
		}

		authHeader := strings.Split(rawAuthHeader, " ")
		if len(authHeader) != 2 && authHeader[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, endpoints.NewErrorResponse("Invalid Authorization header"))
			c.Abort()
			return
		}

		token := authHeader[1]
		claims, err := s.Verify(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, endpoints.ErrorResponseFromError(err))
			c.Abort()
			return
		}

		c.Set(UserIdKey, claims.Uid)
	}
}

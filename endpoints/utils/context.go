package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/services"
)

func GetUserClaimsFromContext(c *gin.Context) *services.UserClaims {
	rawUserClaims, exists := c.Get(middlewares.UserClaimsKey)
	if !exists {
		return nil
	}

	userClaims, ok := rawUserClaims.(*services.UserClaims)
	if !ok {
		return nil
	}

	return userClaims
}

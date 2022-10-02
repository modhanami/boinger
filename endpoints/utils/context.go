package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/services/tokens"
)

func GetUserClaimsFromContext(c *gin.Context) *tokens.UserClaims {
	rawUserClaims, exists := c.Get(middlewares.UserClaimsKey)
	if !exists {
		return nil
	}

	userClaims, ok := rawUserClaims.(*tokens.UserClaims)
	if !ok {
		return nil
	}

	return userClaims
}

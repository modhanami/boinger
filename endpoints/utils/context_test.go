package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/middlewares"
	"github.com/modhanami/boinger/services/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserClaimsFromContext(t *testing.T) {
	tests := []struct {
		name        string
		userClaims  interface{}
		expectedNil bool
	}{
		{
			name:        "return nil when user claims does not exist",
			userClaims:  nil,
			expectedNil: true,
		},
		{
			name:        "return nil when user claims is not user claims",
			userClaims:  "not user claims",
			expectedNil: true,
		},
		{
			name: "return user claims when user claims is user claims",
			userClaims: &tokens.UserClaims{
				ID:       uint(123),
				Username: "user123",
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Set(middlewares.UserClaimsKey, tt.userClaims)

			userClaims := GetUserClaimsFromContext(c)

			if tt.expectedNil {
				assert.Nil(t, userClaims)
			} else {
				assert.NotNil(t, userClaims)
				assert.Equal(t, tt.userClaims, userClaims)
			}
		})
	}
}

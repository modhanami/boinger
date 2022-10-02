package usercontext

import (
	"github.com/modhanami/boinger/services/tokens"
)

type UserContext interface {
	UserID() uint
	Username() string
	Email() string
}

type claimsUserContext struct {
	claims *tokens.UserClaims
}

func (g *claimsUserContext) UserID() uint {
	return g.claims.ID
}

func (g *claimsUserContext) Username() string {
	return g.claims.Username
}

func (g *claimsUserContext) Email() string {
	panic("implement me")
}

var _ UserContext = (*claimsUserContext)(nil)

func NewClaimsUserContext(c *tokens.UserClaims) UserContext {
	return &claimsUserContext{claims: c}
}

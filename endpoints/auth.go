package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/services"
	"net/http"
)

var UserClaimsKey = "userClaims"
var UserIdKey = "userId"
var AuthTokenCookieName = "auth_token"

type UserClaimsResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func NewUserClaimsResponseFromClaims(claims *services.UserClaims) *UserClaimsResponse {
	return &UserClaimsResponse{
		ID:       claims.ID,
		Username: claims.Username,
	}
}

type tokenResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenResponse(token string, refreshToken string) *tokenResponse {
	return &tokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
}

func MakeLoginEndpoint(s services.AuthService, userTokenService services.UserTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		user, err := s.Authenticate(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponseFromError(err))
			return
		}

		token, _, err := userTokenService.Create(&user, services.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
			return
		}

		refreshToken, err := userTokenService.RenewRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
			return
		}

		c.JSON(http.StatusOK, NewTokenResponse(token, refreshToken.Token))
	}
}

func MakeRegisterEndpoint(s services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		_, err := s.Register(username, password)
		if err != nil {
			if err == services.ErrUserAlreadyExists {
				c.JSON(http.StatusConflict, ErrorResponseFromError(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
		}

		c.Status(http.StatusCreated)
	}
}

func MakeUserInfoEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawClaims, exists := c.Get(UserClaimsKey)
		if !exists {
			c.Status(http.StatusInternalServerError)
			return
		}

		claims := rawClaims.(services.UserClaims)
		c.JSON(http.StatusOK, NewUserClaimsResponseFromClaims(&claims))
	}
}

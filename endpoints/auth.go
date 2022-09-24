package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints/response"
	"github.com/modhanami/boinger/endpoints/utils"
	"github.com/modhanami/boinger/services"
	"log"
	"net/http"
	"time"
)

var (
	RefreshTokenCookieName   = "refresh_token"
	RefreshTokenCookieMaxAge = int(30 * 24 * time.Hour / time.Second)
	IsSecureCookieDisabled   bool
)

type UserClaimsResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewLoginResponse(token string) *loginResponse {
	return &loginResponse{
		AccessToken: token,
	}
}

func MakeLoginEndpoint(s services.AuthService, userTokenService services.UserTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		user, err := s.Authenticate(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponseFromError(err))
			return
		}

		token, err := userTokenService.Create(user, services.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponseFromError(err))
			return
		}

		refreshToken, err := userTokenService.RenewRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponseFromError(err))
			return
		}

		log.Println(RefreshTokenCookieMaxAge)
		c.SetCookie(RefreshTokenCookieName, refreshToken.Token, RefreshTokenCookieMaxAge, "/", "", !IsSecureCookieDisabled, true)
		c.JSON(http.StatusOK, NewLoginResponse(token.Value))
	}
}

func MakeRegisterEndpoint(s services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")
		if username == "" || password == "" || email == "" {
			c.JSON(http.StatusBadRequest, response.ErrorResponseFromError(services.ErrInvalidCredentials))
			return
		}

		_, err := s.Register(username, email, password)
		if err != nil {
			if err == services.ErrUserAlreadyExists {
				c.JSON(http.StatusConflict, response.ErrorResponseFromError(err))
				return
			}
			c.JSON(http.StatusInternalServerError, response.ErrorResponseFromError(err))
		}

		c.Status(http.StatusCreated)
	}
}

func MakeUserInfoEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := utils.GetUserClaimsFromContext(c)
		if userClaims == nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusOK, userClaims)
	}
}

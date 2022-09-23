package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/services"
	"log"
	"net/http"
	"time"
)

var (
	UserClaimsKey            = "userClaims"
	RefreshTokenCookieName   = "refresh_token"
	RefreshTokenCookieMaxAge = int(30 * 24 * time.Hour / time.Second)
	IsSecureCookieDisabled   bool
)

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
			c.JSON(http.StatusUnauthorized, ErrorResponseFromError(err))
			return
		}

		token, err := userTokenService.Create(user, services.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
			return
		}

		refreshToken, err := userTokenService.RenewRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
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
			c.JSON(http.StatusBadRequest, ErrorResponseFromError(services.ErrInvalidCredentials))
			return
		}

		_, err := s.Register(username, email, password)
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

		claims := rawClaims.(*services.UserClaims)
		c.JSON(http.StatusOK, NewUserClaimsResponseFromClaims(claims))
	}
}

package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/services"
	"net/http"
	"os"
	"time"
)

var UserClaimsKey = "userClaims"
var UserIdKey = "userId"
var AuthTokenCookieName = "auth_token"

type UserClaimsResponse struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
}

func NewUserClaimsResponseFromClaims(claims *services.UserClaims) *UserClaimsResponse {
	return &UserClaimsResponse{
		Uid:      claims.Uid,
		Username: claims.Username,
	}
}

func MakeLoginEndpoint(s services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		userToken, claims, err := s.Login(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponseFromError(err))
			return
		}

		now := time.Now()
		oneYearFromNow := now.AddDate(1, 0, 0)
		maxAge := int(oneYearFromNow.Sub(now).Seconds())

		disableSecureCookiesEnv := os.Getenv("SECURE_COOKIES_DISABLED")
		secure := disableSecureCookiesEnv != "true"

		c.SetCookie(AuthTokenCookieName, userToken, maxAge, "/", "", secure, true)
		c.JSON(http.StatusOK, NewUserClaimsResponseFromClaims(&claims))
	}
}

func MakeRegisterEndpoint(s services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		success, err := s.Register(username, password)
		if err != nil {
			if err == services.ErrUserAlreadyExists {
				c.JSON(http.StatusConflict, ErrorResponseFromError(err))
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
		}

		if !success {
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

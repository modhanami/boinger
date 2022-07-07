package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"net/http"
	"os"
	"strconv"
	"time"
)

func MakeListEndpoint(s services.BoingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		boings, err := s.List()
		if err != nil {
			c.JSON(500, ErrorResponseFromError(err))
			return
		}
		c.JSON(200, boings)
	}
}

type CreateRequest struct {
	Text string
}

func MakeCreateEndpoint(s services.BoingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(500, ErrorResponseFromError(err))
			return
		}

		err = s.Create(request.Text, 0)
		if err != nil {
			c.JSON(500, ErrorResponseFromError(err))
			return
		}
		c.Status(201)
	}
}

func MakeGetByIdEndpoint(s services.BoingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 0)
		if err != nil {
			c.JSON(500, ErrorResponseFromError(err))
			return
		}

		boing, err := s.Get(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponseFromError(err))
			return
		} else if boing == (models.BoingModel{}) {
			c.JSON(http.StatusNotFound, NewErrorResponse("Boing not found"))
			return
		}

		c.JSON(http.StatusOK, boing)
	}
}

var AuthTokenCookieName = "auth_token"

func MakeLoginEndpoint(s services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		userToken, err := s.Login(username, password)
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
		c.Status(http.StatusOK)
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

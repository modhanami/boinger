package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestVerifyJWTUserTokenMiddleware_Authenticated(t *testing.T) {
	router := makeRouter()
	userTokenService := services.NewUserTokenService()
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	validToken, err := userTokenService.Create(&models.UserModel{
		Uid: "test",
	}, services.CreateOptions{
		Exp: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.FailNow()
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer "+validToken)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_SetUserId(t *testing.T) {
	router := makeRouter()
	userTokenService := services.NewUserTokenService()
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		userId := c.MustGet(UserIdKey).(string)
		assert.Equal(t, "test", userId)
		c.Status(http.StatusOK)
	})

	validToken, err := userTokenService.Create(&models.UserModel{
		Uid: "test",
	}, services.CreateOptions{
		Exp: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.FailNow()
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer "+validToken)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_InvalidAuthHeader(t *testing.T) {
	router := makeRouter()
	userTokenService := services.NewUserTokenService()
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		t.FailNow()
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "Bearer invalid")
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_InvalidAuthScheme(t *testing.T) {
	router := makeRouter()
	userTokenService := services.NewUserTokenService()
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		t.FailNow()
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "invalid")
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_EmptyAuthHeader(t *testing.T) {
	router := makeRouter()
	userTokenService := services.NewUserTokenService()
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

//go:build exclude

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func setup(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fail()
	}

	err = gdb.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		return nil
	}

	return gdb.Debug().Begin()
}

func TestVerifyJWTUserTokenMiddleware_Authenticated(t *testing.T) {
	db := setup(t)
	router := makeRouter()
	userTokenService := services.NewUserTokenService(db)
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	validToken, _, err := userTokenService.Create(&models.User{
		Uid: "test",
	}, services.CreateOptions{
		Exp: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.FailNow()
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  endpoints.AuthTokenCookieName,
		Value: validToken,
	})
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_SetUserId(t *testing.T) {
	db := setup(t)
	router := makeRouter()
	userTokenService := services.NewUserTokenService(db)
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		userId := c.MustGet(endpoints.UserIdKey).(string)
		assert.Equal(t, "test", userId)
		c.Status(http.StatusOK)
	})

	validToken, _, err := userTokenService.Create(&models.User{
		Uid: "test",
	}, services.CreateOptions{
		Exp: time.Now().Add(time.Hour),
	})

	if err != nil {
		t.FailNow()
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  endpoints.AuthTokenCookieName,
		Value: validToken,
	})
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_InvalidAuthCookie(t *testing.T) {
	db := setup(t)
	router := makeRouter()
	userTokenService := services.NewUserTokenService(db)
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		t.FailNow()
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.AddCookie(&http.Cookie{
		Name:  endpoints.AuthTokenCookieName,
		Value: "invalid",
	})
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_EmptyAuthCookie(t *testing.T) {
	db := setup(t)
	router := makeRouter()
	userTokenService := services.NewUserTokenService(db)
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

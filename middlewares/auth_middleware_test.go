package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/modhanami/boinger/endpoints"
	"github.com/modhanami/boinger/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	validToken   = "valid"
	expiredToken = "expired"
	invalidToken = "invalid"
)

//func setup(t *testing.T) *gorm.DB {
//	gdb, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
//	if err != nil {
//		t.Fail()
//	}
//
//	err = gdb.AutoMigrate(&models.User{}, &models.RefreshToken{})
//	if err != nil {
//		return nil
//	}
//
//	return gdb.Debug().Begin()
//}

func TestVerifyJWTUserTokenMiddleware_AuthorizationHeaderPresent(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "authenticated",
			token:          fmt.Sprintf("Bearer %s", validToken),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "expired token",
			token:          fmt.Sprintf("Bearer %s", expiredToken),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			token:          fmt.Sprintf("Bearer %s", invalidToken),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid scheme",
			token:          fmt.Sprintf("Basic %s", validToken),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder, request, router := setupGETRequestHelpers()
			request.Header.Add("Authorization", tt.token)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}

func TestVerifyJWTUserTokenMiddleware_MissingAuthorizationHeader(t *testing.T) {
	recorder, request, router := setupGETRequestHelpers()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestVerifyJWTUserTokenMiddleware_SetUserClaimsInContext(t *testing.T) {
	userTokenService := &fakeUserTokenService{}
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router := makeRouter()
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		raw, exists := c.Get(endpoints.UserClaimsKey)
		assert.True(t, exists)
		claims := raw.(*services.UserClaims)
		assert.NotNil(t, claims)
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", validToken))
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func setupGETRequestHelpers() (*httptest.ResponseRecorder, *http.Request, *gin.Engine) {
	userTokenService := &fakeUserTokenService{}
	userTokenMiddleware := MakeVerifyJWTUserTokenMiddleware(userTokenService)
	router := makeRouter()
	router.GET("/", userTokenMiddleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	return recorder, request, router
}

func makeRouter() *gin.Engine {
	router := gin.Default()
	return router
}

type fakeUserTokenService struct {
	services.UserTokenService
}

func (m *fakeUserTokenService) Verify(token string) (*services.UserClaims, error) {
	switch token {
	case "valid":
		return &services.UserClaims{}, nil
	case "expired":
		return nil, fmt.Errorf("token is expired")
	default:
		return nil, fmt.Errorf("token is invalid")
	}
}

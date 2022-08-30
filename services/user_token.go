package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/modhanami/boinger/models"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var (
	ErrUserClaimsParseFailed      = errors.New("failed to parse user claims")
	ErrInvalidToken               = errors.New("invalid token")
	ErrFailedToCreateRefreshToken = errors.New("failed to create refresh token")
)

type UserTokenService interface {
	Create(user *models.User, options CreateOptions) (UserToken, UserClaims, error)
	Verify(token string) (UserClaims, error)
	RenewRefreshToken(userId uint) (*models.RefreshToken, error)
}

type userTokenService struct {
	db *gorm.DB
}

func NewUserTokenService(db *gorm.DB) UserTokenService {
	return &userTokenService{
		db: db,
	}
}

var hmacSampleSecret = []byte("dont-mind-me-this-is-a-secret")

type CreateOptions struct {
	Exp time.Time
}

type UserToken = string

func (s *userTokenService) Create(user *models.User, options CreateOptions) (UserToken, UserClaims, error) {
	var exp time.Time
	if options.Exp.IsZero() {
		exp = time.Now().Add(time.Hour * 24 * 7)
	} else {
		exp = options.Exp
	}

	claims := NewUserClaimsWithExp(user.Uid, user.Username, exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, claims, err
}

type UserClaims struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewUserClaims(uid string, username string) UserClaims {
	oneWeekFromNow := time.Now().AddDate(0, 0, 7)
	return NewUserClaimsWithExp(uid, username, oneWeekFromNow)
}

func NewUserClaimsWithExp(uid string, username string, exp time.Time) UserClaims {
	return UserClaims{
		Uid:      uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "boinger",
		},
	}
}

func (s *userTokenService) Verify(rawToken string) (UserClaims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return UserClaims{}, ErrUserClaimsParseFailed
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return UserClaims{}, ErrInvalidToken
	}

	return NewUserClaims(claims.Uid, claims.Username), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-"

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *userTokenService) RenewRefreshToken(userId uint) (*models.RefreshToken, error) {
	refreshTokenString := GenerateRandomString(64)
	refreshToken := models.NewRefreshToken(userId, refreshTokenString)

	if err := s.db.Model(&models.RefreshToken{}).Where("user_id = ?", userId).Update("revoked_at", time.Now()).Error; err != nil {
		return nil, ErrFailedToCreateRefreshToken
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return nil, ErrFailedToCreateRefreshToken
	}

	return refreshToken, nil
}

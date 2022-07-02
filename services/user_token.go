package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/modhanami/boinger/models"
	"time"
)

type UserTokenService interface {
	Create(user *models.UserModel, options CreateOptions) (UserToken, error)
	Verify(token string) (Claims, error)
}

type userTokenService struct{}

func NewUserTokenService() UserTokenService {
	return &userTokenService{}
}

var hmacSampleSecret = []byte("dont-mind-me-this-is-a-secret")

type CreateOptions struct {
	Exp time.Time
}

type UserToken = string

func (s *userTokenService) Create(user *models.UserModel, options CreateOptions) (UserToken, error) {
	var exp int64
	if options.Exp.IsZero() {
		exp = time.Now().Add(time.Hour * 24 * 7).Unix()
	} else {
		exp = options.Exp.Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(user.Uid, exp).ToJWTClaims())

	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

type Claims struct {
	Uid string
	Exp int64
}

func NewClaims(uid string, exp int64) Claims {
	return Claims{
		Uid: uid,
		Exp: exp,
	}
}

func (c Claims) ToJWTClaims() jwt.Claims {
	return jwt.MapClaims{
		"uid": c.Uid,
		"exp": c.Exp,
	}
}

func (s *userTokenService) Verify(rawToken string) (Claims, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	if _, ok := claims["uid"]; !ok {
		return Claims{}, errors.New("invalid token")
	}

	uidStr, ok := claims["uid"].(string)
	if !ok {
		return Claims{}, errors.New("invalid token")
	}

	expInt64, ok := claims["exp"].(float64)
	if !ok {
		return Claims{}, errors.New("invalid token")
	}

	return Claims{
		Uid: uidStr,
		Exp: int64(expInt64),
	}, nil
}

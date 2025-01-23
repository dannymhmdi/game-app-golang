package authservice

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"mymodule/entity"
	"strings"
	"time"
)

type Config struct {
	SigningKey             string
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	RefreshSubject         string
	AccessSubject          string
}

type Service struct {
	config Config
}

type CustomClaims struct {
	UserId           uint
	Name             string
	RegisteredClaims jwt.RegisteredClaims
}

func New(cfg Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.CreateToken(user, s.config.AccessTokenExpireTime, "at")
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.CreateToken(user, s.config.RefreshTokenExpireTime, "rt")
}

func (s Service) CreateToken(user entity.User, expireTime time.Duration, tokenSubject string) (string, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(expireTime))
	key := []byte(s.config.SigningKey)
	claims := CustomClaims{
		UserId: user.ID,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			Subject:   tokenSubject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func (s Service) ParseToken(tokenString string) (*CustomClaims, error) { // Your secret key
	//parsedToken:=&CustomClaims{}
	tokenString = strings.Split(tokenString, " ")[1]
	//tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	fmt.Println("tokenString", tokenString)
	key := []byte(s.config.SigningKey)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	customClaims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return customClaims, nil
}

func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.ExpiresAt, nil
}

func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.IssuedAt, nil
}

func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.NotBefore, nil
}

func (c CustomClaims) GetIssuer() (string, error) {
	return c.RegisteredClaims.Issuer, nil
}

func (c CustomClaims) GetSubject() (string, error) {
	return c.RegisteredClaims.Subject, nil
}

func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	return c.RegisteredClaims.Audience, nil
}

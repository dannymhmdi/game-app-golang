package authService

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"mymodule/entity"
	"strings"
	"time"
)

type Config struct {
	SigningKey             string        `koanf:"signKey"`
	AccessTokenExpireTime  time.Duration `koanf:"access_token_expire_time"`
	RefreshTokenExpireTime time.Duration `koanf:"refresh_token_expire_time"`
	RefreshSubject         string        `koanf:"refresh_subject"`
	AccessSubject          string        `koanf:"access_subject"`
}

type Service struct {
	config Config
	repo   AuthRepositoryService
}

type AuthRepositoryService interface {
	StoreRefreshToken(ctx context.Context, refreshToken string, userInfo entity.User, cfg Config) error
}

type CustomClaims struct {
	UserId           uint
	Name             string
	Role             entity.Role
	RegisteredClaims jwt.RegisteredClaims
}

func New(cfg Config, repo AuthRepositoryService) *Service {
	return &Service{
		config: cfg,
		repo:   repo,
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
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			Subject:   tokenSubject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func (s Service) ParseToken(tokenString string) (*CustomClaims, error) {
	if isRefreshToken := strings.Contains(tokenString, "refresh-token"); isRefreshToken {
		tokenString = strings.Replace(tokenString, "refresh-token=", "", 1)
	} else {
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	}

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

func (s Service) StoreToken(ctx context.Context, refreshToken string, userInfo entity.User, cfg Config, userAgent string) error {

	sErr := s.repo.StoreRefreshToken(ctx, refreshToken, userInfo, cfg)
	if sErr != nil {
		return sErr
	}

	return nil
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

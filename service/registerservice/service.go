package registerservice

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
	"mymodule/pkg/validation/passwordvalidation"
	"mymodule/pkg/validation/phonenumbervalidation"
	"time"
)

type RegisterRepositoryService interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	RegisterUser(user entity.User) (entity.User, error)
	GetUserById(id uint) (entity.User, error)
}

type Service struct {
	repository RegisterRepositoryService
}

func New(rep RegisterRepositoryService) *Service {
	return &Service{
		rep,
	}
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
	Password    string
}

type RegisterResponse struct {
	User entity.User
}

type LoginRequest struct {
	PhoneNumber string
	Password    string
}

type LoginResponse struct {
	Message     string
	Status      bool
	user        entity.User
	AccessToken string
}

type ProfileRequest struct {
	Id uint
}

type ProfileResponse struct {
	User entity.User
}

type CustomClaims struct {
	UserId           uint
	Name             string
	RegisteredClaims jwt.RegisteredClaims
}

func (s Service) RegisterUser(req RegisterRequest) (RegisterResponse, error) {
	//TODO - hashing password
	fmt.Printf("RegisterRequest:%+v\n", req)
	if iErr := phonenumbervalidation.IsValid(req.PhoneNumber); iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	if isUnique, iErr := s.repository.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	if iErr := passwordvalidation.IsPasswordValid(req.Password); iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if rErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", rErr)
	}
	return RegisterResponse{User: createdUser}, nil

}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, gErr := s.repository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {
		return LoginResponse{}, gErr
	}
	fmt.Println("uuu", user)
	if cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); cErr != nil {
		return LoginResponse{}, fmt.Errorf("password is incorrect :%v\n", cErr)
	}

	accessToken, gErr := GenerateJWT(user.ID, user.Name)
	if gErr != nil {
		return LoginResponse{}, fmt.Errorf("failed to generate token: %v\n", gErr)
	}

	token, vErr := validateJWT(accessToken)
	if vErr != nil {
		return LoginResponse{}, fmt.Errorf("failed to validate token: %v\n", vErr)
	}
	fmt.Println("token", token)
	return LoginResponse{Message: "success", Status: true, user: user, AccessToken: accessToken}, nil
}

func (s Service) GetUserProfile(req ProfileRequest) (ProfileResponse, error) {
	userInfo, gErr := s.repository.GetUserById(req.Id)
	if gErr != nil {
		return ProfileResponse{}, gErr
	}

	return ProfileResponse{User: userInfo}, nil
}

func GenerateJWT(id uint, name string) (string, error) {
	mySigningKey := []byte("secret") // Your secret key

	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 1))

	claims := CustomClaims{
		UserId: id,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func validateJWT(tokenString string) (*CustomClaims, error) {
	mySigningKey := []byte("secret") // Your secret key

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(*CustomClaims), nil
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

package registerservice

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"mymodule/pkg/validation/passwordvalidation"
	"mymodule/pkg/validation/phonenumbervalidation"
)

type RegisterRepositoryService interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	RegisterUser(user entity.User) (entity.User, error)
	GetUserById(id uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	authSvc    AuthGenerator
	repository RegisterRepositoryService
}

func New(rep RegisterRepositoryService, authsvc AuthGenerator) *Service {
	return &Service{
		authSvc:    authsvc,
		repository: rep,
	}
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User `json:"user"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	user    entity.User `json:"user"`
	Token   Token       `json:"token"`
}

type ProfileRequest struct {
	Id uint `json:"id"`
}

type ProfileResponse struct {
	User entity.User `json:"user"`
}

func (s Service) RegisterUser(req RegisterRequest) (RegisterResponse, error) {
	//TODO - hashing password
	fmt.Printf("RegisterRequest:%+v\n", req)
	if iErr := phonenumbervalidation.IsValid(req.PhoneNumber); iErr != nil {
		return RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(iErr).
			SetMsg("password is not valid").SetKind(richerr.KindInvalid)
	}

	if isUnique, iErr := s.repository.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || iErr != nil {
		return RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(iErr).
			SetMsg("phone number is not unique").
			SetKind(richerr.KindInvalid)
	}

	if iErr := passwordvalidation.IsPasswordValid(req.Password); iErr != nil {
		return RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(iErr).
			SetMsg("password is not valid").
			SetKind(richerr.KindInvalid)
	}

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})

	if rErr != nil {
		return RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(rErr).
			SetMsg("failed to register user").
			SetKind(richerr.KindInvalid)
	}

	return RegisterResponse{User: createdUser}, nil

}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, gErr := s.repository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {
		return LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetWrappedErr(gErr)
	}

	if cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); cErr != nil {
		return LoginResponse{}, fmt.Errorf("password is incorrect :%v\n", cErr)
	}

	accessToken, gErr := s.authSvc.CreateAccessToken(user)
	if gErr != nil {
		return LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetMsg("failed to generate access-token").
			SetWrappedErr(gErr)
	}

	refreshToken, gErr := s.authSvc.CreateRefreshToken(user)
	if gErr != nil {
		return LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetWrappedErr(gErr).
			SetMsg("failed to generate refresh-token")
	}

	return LoginResponse{Message: "success", Status: true, user: user, Token: Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}, nil
}

func (s Service) GetUserProfile(id uint) (ProfileResponse, error) {
	userInfo, gErr := s.repository.GetUserById(id)

	if gErr != nil {
		return ProfileResponse{}, richerr.New().
			SetOperation("registerService.GetUserProfile").
			SetMsg("cant find user").
			SetWrappedErr(gErr).
			SetKind(richerr.KindNotFound)
	}

	return ProfileResponse{User: userInfo}, nil
}

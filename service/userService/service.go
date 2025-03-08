package userService

import (
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/validator/uservalidator"
)

type RegisterRepositoryService interface {
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	RegisterUser(user entity.User) (entity.User, error)
	GetUserById(id uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	validator  uservalidator.Validator
	authSvc    AuthGenerator
	repository RegisterRepositoryService
}

func New(rep RegisterRepositoryService, authsvc AuthGenerator, validator uservalidator.Validator) *Service {
	return &Service{
		validator:  validator,
		authSvc:    authsvc,
		repository: rep,
	}
}

func (s Service) Register(req params.RegisterRequest) (params.RegisterResponse, error) {

	//TODO - wirte validator methods here

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Role:        entity.UserRole,
	})

	if rErr != nil {
		return params.RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(rErr).
			SetMsg("failed to register user").
			SetKind(richerr.KindInvalid)
	}

	return params.RegisterResponse{User: createdUser}, nil
}

func (s Service) Login(req params.LoginRequest) (params.LoginResponse, error) {

	user, vErr := s.repository.GetUserByPhoneNumber(req.PhoneNumber)
	if vErr != nil {
		return params.LoginResponse{}, vErr
	}

	accessToken, gErr := s.authSvc.CreateAccessToken(user)
	if gErr != nil {
		return params.LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetMsg("failed to generate access-token").
			SetWrappedErr(gErr)
	}

	refreshToken, gErr := s.authSvc.CreateRefreshToken(user)
	if gErr != nil {
		return params.LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetWrappedErr(gErr).
			SetMsg("failed to generate refresh-token")
	}

	return params.LoginResponse{Message: "success", Status: true, Token: params.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}, nil
}

func (s Service) GetUserProfile(id uint) (params.ProfileResponse, error) {
	userInfo, gErr := s.repository.GetUserById(id)

	if gErr != nil {
		return params.ProfileResponse{}, richerr.New().
			SetOperation("registerService.GetUserProfile").
			SetMsg("cant find user").
			SetWrappedErr(gErr).
			SetKind(richerr.KindNotFound)
	}

	return params.ProfileResponse{User: userInfo}, nil
}

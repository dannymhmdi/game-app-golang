package registerservice

import (
	"mymodule/dto"
	"mymodule/entity"
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

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {

	//TODO - wirte validator methods here
	isValid, vErr := s.validator.ValidateRegisterCredentials(req)
	if vErr != nil || !isValid {
		return dto.RegisterResponse{}, richerr.New().
			SetMsg(vErr.Error()).
			SetOperation("reigsterservice.RegisterUser").
			SetKind(richerr.KindInvalid)
	}

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})

	if rErr != nil {
		return dto.RegisterResponse{}, richerr.New().
			SetOperation("registerService.RegisterUser").
			SetWrappedErr(rErr).
			SetMsg("failed to register user").
			SetKind(richerr.KindInvalid)
	}

	return dto.RegisterResponse{User: createdUser}, nil
}

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {

	user, vErr := s.validator.ValidateLoginCredentials(req)
	if vErr != nil {
		return dto.LoginResponse{}, richerr.New().
			SetMsg(vErr.Error()).
			SetOperation("registerService.Login").
			SetKind(richerr.KindInvalid)
	}

	accessToken, gErr := s.authSvc.CreateAccessToken(user)
	if gErr != nil {
		return dto.LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetMsg("failed to generate access-token").
			SetWrappedErr(gErr)
	}

	refreshToken, gErr := s.authSvc.CreateRefreshToken(user)
	if gErr != nil {
		return dto.LoginResponse{}, richerr.New().
			SetOperation("registerService.Login").
			SetWrappedErr(gErr).
			SetMsg("failed to generate refresh-token")
	}

	return dto.LoginResponse{Message: "success", Status: true, Token: dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}, nil
}

func (s Service) GetUserProfile(id uint) (dto.ProfileResponse, error) {
	userInfo, gErr := s.repository.GetUserById(id)

	if gErr != nil {
		return dto.ProfileResponse{}, richerr.New().
			SetOperation("registerService.GetUserProfile").
			SetMsg("cant find user").
			SetWrappedErr(gErr).
			SetKind(richerr.KindNotFound)
	}

	return dto.ProfileResponse{User: userInfo}, nil
}

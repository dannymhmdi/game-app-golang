package registerservice

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
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

	accessToken, gErr := s.authSvc.CreateAccessToken(user)
	if gErr != nil {
		return LoginResponse{}, fmt.Errorf("failed to generate access-token: %v\n", gErr)
	}

	refreshToken, gErr := s.authSvc.CreateRefreshToken(user)
	if gErr != nil {
		return LoginResponse{}, fmt.Errorf("failed to generate refresh-token: %v\n", gErr)
	}

	return LoginResponse{Message: "success", Status: true, user: user, Token: Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}}, nil
}

func (s Service) GetUserProfile(id uint) (ProfileResponse, error) {
	userInfo, gErr := s.repository.GetUserById(id)
	if gErr != nil {
		return ProfileResponse{}, gErr
	}

	return ProfileResponse{User: userInfo}, nil
}

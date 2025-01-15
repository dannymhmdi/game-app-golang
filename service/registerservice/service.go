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
	Message string
	Status  bool
	user    entity.User
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
	fmt.Println("userff", user)

	if cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); cErr != nil {
		return LoginResponse{}, fmt.Errorf("password is incorrect :%v\n", cErr)
	}

	return LoginResponse{Message: "success", Status: true, user: user}, nil
}

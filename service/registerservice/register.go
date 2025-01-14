package registerservice

import (
	"fmt"
	"mymodule/entity"
	"mymodule/pkg/phonenumbervalidation"
)

type RegisterRepositoryService interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
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
}

type RegisterResponse struct {
	User entity.User
}

func (s Service) RegisterUser(req RegisterRequest) (RegisterResponse, error) {
	fmt.Printf("RegisterRequest:%+v\n", req)
	if iErr := phonenumbervalidation.IsValid(req.PhoneNumber); iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	if isUnique, iErr := s.repository.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	})
	if rErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", rErr)
	}
	return RegisterResponse{User: createdUser}, nil

}

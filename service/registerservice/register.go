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

type RegisterRequest struct {
	Name        string
	phoneNumber string
}

type RegisterResponse struct {
	User entity.User
}

func (s Service) RegisterUser(req RegisterRequest) (RegisterResponse, error) {
	if isValid, iErr := phonenumbervalidation.IsValid(req.phoneNumber); iErr != nil || !isValid {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	if isUnique, iErr := s.repository.IsPhoneNumberUnique(req.phoneNumber); !isUnique || iErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", iErr)
	}

	createdUser, rErr := s.repository.RegisterUser(entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.phoneNumber,
	})
	if rErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error: %v\n", rErr)
	}
	return RegisterResponse{User: createdUser}, nil

}

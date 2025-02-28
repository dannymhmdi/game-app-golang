package uservalidator

import "mymodule/entity"

type ValidatorRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repository ValidatorRepository
}

func New(repo ValidatorRepository) *Validator {
	return &Validator{repository: repo}
}

package uservalidator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"mymodule/dto"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"regexp"
	"unicode"
)

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

func (v Validator) ValidateRegisterCredentials(credential dto.RegisterRequest) error {
	vErr := validation.ValidateStruct(&credential,
		validation.Field(&credential.Name, validation.Required, validation.Length(4, 50)),

		validation.Field(&credential.PhoneNumber,
			validation.Required, validation.Length(11, 11),
			validation.Match(regexp.MustCompile(`^09\d{9}$`)),
			validation.By(v.IsPhoneNumberUnique)),

		validation.Field(&credential.Password,
			validation.Required,
			validation.Length(8, 100),
			validation.By(v.PasswordValidation),
		),
	)

	if vErr != nil {
		return richerr.New().
			SetMsg(vErr.Error()).
			SetWrappedErr(vErr).
			SetOperation("uservalidator.ValidateCredentials").
			SetKind(richerr.KindInvalid)
	}

	return nil
}

func (v Validator) IsPhoneNumberUnique(value interface{}) error {
	phoneNumber := value.(string)

	isUnique, iErr := v.repository.IsPhoneNumberUnique(phoneNumber)
	if iErr != nil {
		return richerr.New().
			SetOperation("uservalidator.IsPhoneNumberUnique").
			SetMsg(iErr.Error()).
			SetKind(richerr.KindInvalid)
	}

	if !isUnique {
		return richerr.New().
			SetOperation("uservalidator.IsPhoneNumberUnique").
			SetMsg("phone number is not unique").
			SetKind(richerr.KindInvalid)
	}
	return nil
}

func (v Validator) PasswordValidation(value interface{}) error {
	password := value.(string)

	var hasDigit, hasUpperCase, hasSpecialChar bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	if !hasDigit || !hasUpperCase || !hasSpecialChar {
		var err string
		if !hasDigit {
			err += "password must contain at least one digit\n"
		}

		if !hasUpperCase {
			err += "password must contain at least one upper case\n"
		}

		if !hasSpecialChar {
			err += "password must contain at least one special character\n"
		}

		return richerr.New().
			SetMsg(err).
			SetOperation("uservalidator.PasswordValidation").
			SetKind(richerr.KindInvalid)
	}

	return nil
}

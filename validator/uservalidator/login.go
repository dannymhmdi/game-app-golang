package uservalidator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"mymodule/dto"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"regexp"
)

func (v Validator) ValidateLoginCredentials(credential dto.LoginRequest) (entity.User, error) {
	err := validation.ValidateStruct(&credential,
		validation.Field(&credential.PhoneNumber,
			validation.Required, validation.Length(11, 11),
			validation.Match(regexp.MustCompile(`^09\d{9}$`))))

	if err != nil {
		return entity.User{}, err
	}

	user, gErr := v.repository.GetUserByPhoneNumber(credential.PhoneNumber)
	if gErr != nil {
		return entity.User{}, gErr
	}

	if cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password)); cErr != nil {
		return entity.User{}, richerr.New().
			SetMsg("password is incorrect").
			SetOperation("userValidator.ValidateLoginCredentials").
			SetKind(richerr.KindInvalid)
	}

	return user, nil
}

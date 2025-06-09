package authValidator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"mymodule/pkg/richerr"
)

type Validator struct {
	repo              ValidatorRepository
	cookieTokenUserId uint
}

type ValidatorRepository interface {
	IsRefreshTokenValid(cookieRefreshToken string, userId uint) error
}

func New(repo ValidatorRepository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v *Validator) ValidateRefreshToken(cookieRefreshToken string, userid uint) error {
	v.cookieTokenUserId = userid
	vErr := validation.Validate(cookieRefreshToken, validation.Required, validation.By(v.RefreshTokenValidator))
	if vErr != nil {
		return richerr.New().
			SetMsg(vErr.Error()).
			SetKind(richerr.KindUnauthorized).
			SetOperation("authValidator.ValidateRefreshToken").
			SetWrappedErr(vErr)
	}

	return nil

}

func (v *Validator) RefreshTokenValidator(value interface{}) error {
	cookieRefreshToken := value.(string)
	iErr := v.repo.IsRefreshTokenValid(cookieRefreshToken, v.cookieTokenUserId)
	if iErr != nil {
		return richerr.New().
			SetMsg(iErr.Error()).
			SetKind(richerr.KindUnauthorized).
			SetOperation("authValidator.RefreshTokenValidator").
			SetWrappedErr(iErr)
	}

	return nil
}

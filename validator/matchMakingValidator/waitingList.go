package matchMakingValidator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"mymodule/entity"
	"mymodule/params"
	"mymodule/pkg/richerr"
)

type Validator struct {
}

func New() *Validator {
	return &Validator{}
}

func (v Validator) ValidateMatchMakingCredentials(credential params.AddToWaitingListRequest) error {
	err := validation.ValidateStruct(&credential,
		validation.Field(&credential.Category,
			validation.Required,
			validation.By(v.IsCategoryValid)))

	if err != nil {
		return richerr.New().
			SetMsg(err.Error()).
			SetOperation("matchMakingValidator.ValidateMatchMakingCredentials").
			SetKind(richerr.KindInvalid)
	}
	return nil
}

func (v Validator) IsCategoryValid(value interface{}) error {
	category, ok := value.(entity.Category)
	if ok {
		if category.IsValid() {
			return nil
		} else {
			return richerr.New().
				SetMsg("category is not valid").
				SetKind(richerr.KindInvalid).
				SetOperation("matchMakingValidator.IsCategoryValid")
		}
	} else {
		return richerr.New().
			SetMsg("category type not valid").
			SetOperation("matchMakingValidator.IsCategoryValid").
			SetKind(richerr.KindUnexpected)
	}

}

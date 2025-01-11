package phonenumbervalidation

import (
	"fmt"
	"regexp"
)

func IsValid(phonenumber string) (bool, error) {
	regex := `^0\d{10}$`
	match, mErr := regexp.MatchString(regex, phonenumber)
	if mErr != nil {
		return false, fmt.Errorf("failed to check phonenumber regex:%v\n", mErr.Error())
	}

	if !match {
		return false, fmt.Errorf("phonenumber is not valid")
	}

	return true, nil

}

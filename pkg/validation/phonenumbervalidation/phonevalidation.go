package phonenumbervalidation

import (
	"fmt"
	"regexp"
)

func IsValid(phonenumber string) error {
	regex := `^0\d{10}$`
	res, cErr := regexp.Compile(regex)
	if cErr != nil {
		return cErr
	}
	if !res.MatchString(phonenumber) {
		return fmt.Errorf("phonenumber: %s is not valid", phonenumber)
	}
	return nil
	//match, mErr := regexp.MatchString(regex, phonenumber)
	//if mErr != nil {
	//	return fmt.Errorf("failed to check phonenumber regex:%v\n", mErr.Error())
	//}
	//
	//if !match {
	//	return fmt.Errorf("phonenumber is not valid")
	//}
	//
	//return nil

}

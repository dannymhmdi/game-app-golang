package mysql

import "mymodule/entity"

type DB struct {
}

func (d DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {

}

func (d DB) RegisterUser(user entity.User) (entity.User, error) {

}

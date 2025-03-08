package backofficeService

import "mymodule/entity"

type Service struct {
}

func New() *Service {
	return &Service{}
}
func (s Service) ListUsers() ([]entity.User, error) {
	return []entity.User{
		{Name: "fake", ID: 0, PhoneNumber: "fake", Role: entity.Admin, Password: "fake"},
	}, nil
}

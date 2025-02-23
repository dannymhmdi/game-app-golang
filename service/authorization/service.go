package authorization

import (
	"mymodule/entity"
	"mymodule/pkg/richerr"
)

type AuthorizationRepositoryService interface {
	GetUserAcl(userid uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repository AuthorizationRepositoryService
}

func New(repo AuthorizationRepositoryService) *Service {
	return &Service{
		repository: repo,
	}
}

func (s Service) CheckAccess(userId uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {
	userPermissons, gErr := s.repository.GetUserAcl(userId, role)
	if gErr != nil {
		return false, richerr.New().
			SetOperation("authorization.CheckAccess").
			SetWrappedErr(gErr).
			SetMsg(gErr.Error())
	}

	for _, value := range userPermissons {
		for _, permission := range permissions {
			if permission == value {
				return true, nil
			}
		}
	}

	return false, nil
}

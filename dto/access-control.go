package dto

import "mymodule/entity"

type PermissionRequest struct {
	PermissionTitles []entity.PermissionTitle `json:"permission_titles"`
}

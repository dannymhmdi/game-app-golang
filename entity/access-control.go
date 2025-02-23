package entity

type AccessControl struct {
	Id           uint
	ActorId      uint
	ActorType    ActorType
	PermissionId uint
}

type ActorType string

const (
	RoleActorType = "role"
	UserActorType = "user"
)

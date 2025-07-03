package entity

type Role uint

const (
	UserRole Role = iota + 1
	Admin
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return "user"
	case Admin:
		return "admin"
	default:
		return "unknown"
	}
}

func (r *Role) RoleId(role string) {
	switch role {
	case "user":
		*r = UserRole
	case "admin":
		*r = Admin
	}
}

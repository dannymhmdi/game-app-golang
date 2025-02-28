package entity

type Category string

const (
	SoccorCategory Category = "soccer"
)

func (c Category) IsValid() bool {
	switch c {
	case SoccorCategory:
		return true
	default:
		return false
	}
}

package entity

type Category string

const (
	SoccorCategory  Category = "soccer"
	HistoryCategory          = "history"
)

func (c Category) IsValid() bool {
	switch c {
	case SoccorCategory:
		return true
	case HistoryCategory:
		return true
	default:
		return false
	}
}

package entity

type WaitingMember struct {
	UserID    uint
	Timestamp int64
	Category  Category
}

type MatchedPlayers struct {
	UserIDs   []uint
	Category  Category
	Timestamp int64
}

package entity

type RefreshToken struct {
	Id         uint
	UserId     uint
	TokenHash  string
	ExpireAt   []uint8
	Revoked    bool
	DeviceInfo string
	CreatedAt  []uint8
}

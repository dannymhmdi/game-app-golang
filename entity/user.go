package entity

type User struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"-"`
	Role        Role   `json:"role"`
}

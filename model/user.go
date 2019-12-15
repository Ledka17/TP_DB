package model

type User struct {
	Id			int32	`db:"id" json:"-"`
	About		string	`db:"about" json:"about"`
	Email		string	`db:"email" json:"email"`
	Fullname	string	`db:"fullname" json:"fullname"`
	Nickname	string	`db:"nickname" json:"nickname"`
}

type UserUpdate struct {
	About		string
	Email		string
	Fullname	string
}
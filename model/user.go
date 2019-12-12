package model

type User struct {
	About		string
	Email		string
	Fullname	string
	Nickname	string
}

type UserUpdate struct {
	About		string
	Email		string
	Fullname	string
}
package model

type Thread struct {
	Author		string		`db:"author" json:"author"`
	UserId		int32		`db:"user_id" json:"-"`
	Created		string		`db:"created" json:"created"`
	Forum		string		`db:"forum" json:"forum"`
	Id			int32		`db:"id" json:"id"`
	Message		string		`db:"message" json:"message"`
	Slug		string		`db:"slug" json:"slug"`
	Title		string		`db:"title" json:"title"`
	Votes		int32		`db:"votes" json:"votes"`
}

type ThreadUpdate struct {
	Message	string			`json:"message"`
	Title	string			`json:"title"`
}

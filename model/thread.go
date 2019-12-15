package model

import "time"

type Thread struct {
	Author		string		`db:"-" json:"author"`
	UserId		int32		`db:"user_id" json:"-"`
	Created		time.Time	`db:"created" json:"created"`
	Forum		string		`db:"-" json:"forum"`
	ForumId		int32		`db:"forum_id" json:"-"`
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

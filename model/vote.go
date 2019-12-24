package model

type Vote struct {
	Id 			int32	`db:"id" json:"-"`
	ThreadId	int32	`db:"thread_id" json:"-"`
	Nickname	string	`db:"nickname" json:"nickname"`
	Voice		int32	`db:"voice" json:"voice"`
}

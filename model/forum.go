package model

type Forum struct {
	Id		int32	`db:"id" json:"-"`
	Posts	int64	`db:"posts" json:"posts"`
	Slug	string	`db:"slug" json:"slug"`
	Treads	int32	`db:"threads" json:"threads"`
	Title	string	`db:"title" json:"title"`
	User	string	`db:"-" json:"user"`
	UserId	int32	`db:"user_id" json:"-"`
}

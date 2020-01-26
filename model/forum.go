package model

type Forum struct {
	Posts	int64	`db:"posts" json:"posts"`
	Slug	string	`db:"slug" json:"slug"`
	Treads	int32	`db:"threads" json:"threads"`
	Title	string	`db:"title" json:"title"`
	User	string	`db:"author" json:"user"`
}

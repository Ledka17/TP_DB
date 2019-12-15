package model

type Post struct {
	Id			int64	`db:"id" json:"id"`
	Author		string	`db:"-" json:"author"`
	UserId		string	`db:"user_id" json:"-"`
	Created		string	`db:"created" json:"created"`
	Forum		string	`db:"-" json:"forum"`
	ForumId		int32	`db:"forum_id" json:"-"`
	IsEdited	bool	`db:"isEdited" json:"isEdited"`
	Message		string	`db:"message" json:"message"`
	Parent		int64	`db:"parent" json:"parent"`
	Tread		string	`db:"-" json:"thread"`
	ThreadId	int32	`db:"thread_id" json:"-"`
}

type PostFull struct {
	Author	User
	Forum	Forum
	Post 	Post
	Thread	Thread
}

type PostUpdate struct {
	Message	string	`db:"message" json:"message"`
}

package model

type Post struct {
	Id			int64		`db:"id" json:"id"`
	Author		string		`db:"-" json:"author"`
	UserId		int32		`db:"user_id" json:"-"`
	Created		string		`db:"created" json:"created"`
	Forum		string		`db:"-" json:"forum"`
	ForumId		int32		`db:"forum_id" json:"-"`
	IsEdited	bool		`db:"isedited" json:"isEdited"`
	Message		string		`db:"message" json:"message"`
	Parent		int64		`db:"parent" json:"parent"`
	ThreadId	int32		`db:"thread_id" json:"thread"`
}

type PostFull struct {
	Author	*User			`json:"author,omitempty"`
	Forum	*Forum			`json:"forum,omitempty"`
	Post 	*Post			`json:"post,omitempty"`
	Thread	*Thread			`json:"thread,omitempty"`
}

type PostUpdate struct {
	Message	string			`json:"message"`
}

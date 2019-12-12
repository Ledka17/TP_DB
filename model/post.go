package model

type Post struct {
	Author		string
	Created		string
	Forum		string
	Id			int64
	IsEdited	bool
	Message		string
	Parent		int64
	Tread		int32
}

type PostFull struct {
	Author	User
	Forum	Forum
	Post 	Post
	Thread	Thread
}

type PostUpdate struct {
	Message	string
}

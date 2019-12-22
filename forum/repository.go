package forum

import "github.com/Ledka17/TP_DB/model"

type Repository interface {
	IsForumInDB(slug string) bool
	GetForumInDB(slug string) model.Forum
	CreateForumInDB(forum model.Forum) model.Forum
	GetForumUsersInDB(slug string, limit int, since string, desc bool) []model.User

	IsPostInDB(id int) bool
	GetPostInDB(id int, related []string) model.PostFull
	GetPostsInDB(threadSlugOrId string, limit int, since int, sort string, desc bool) []model.Post
	ChangePostInDB(id int, update model.PostUpdate) model.Post
	CreatePostsInDB(posts []model.Post) []model.Post

	CleanUp()
	GetStatusDB() model.Status

	IsUserInDB(nickname string, email string) bool
	GetUserInDB(nickname string, email string) model.User
	GetUsersInDB(nickname string, email string) []model.User
	СreateUserInDB(nickname string, user model.User) model.User
	GetUserIdByName(nickname string) int32

	IsThreadInDB(slugOrId string) bool
	GetThreadInDB(slugOrID string) model.Thread
	CreateThreadInDB(forumSlug string, thread model.Thread) model.Thread
	GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread
	CheckParentPost(posts []model.Post) bool
	ChangeThreadInDB(threadUpdate model.ThreadUpdate, slugOrId string) model.Thread

	VoteForThreadInDB(threadSlugOrId string, vote model.Vote) model.Thread
}

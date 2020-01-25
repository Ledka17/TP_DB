package forum

import "github.com/Ledka17/TP_DB/model"

type Repository interface {
	IsForumInDB(slug string) bool
	GetForumById(id int32) model.Forum
	GetForumInDB(slug string) model.Forum
	CreateForumInDB(forum model.Forum) model.Forum
	GetForumUsersInDB(slug string, limit int, since string, desc bool) []model.User

	IsPostInDB(id int) bool
	GetPostInDB(id int) model.Post
	GetPostsInDB(threadSlugOrId string, limit int, since int, sort string, desc bool) []model.Post
	ChangePostInDB(id int, update model.PostUpdate) (model.Post, error)
	CreatePostsInDB(posts []model.Post, thread model.Thread) []model.Post

	CleanUp()
	GetStatusDB() model.Status

	IsUserInDB(nickname string, email string) bool
	IsUsersInDB(nickname []string) bool
	ChangeUserInDB(nickname string, userUpdate model.UserUpdate) model.User
	GetUserInDB(nickname string, args ...string) model.User
	GetUsersInDB(nickname string, email string) []*model.User
	CreateUserInDB(nickname string, user model.User) model.User
	GetUserIdByName(nickname string) int32

	IsThreadInDB(slugOrId string) bool
	GetThreadInDB(slugOrID string) model.Thread
	GetThreadById(id int) model.Thread
	CreateThreadInDB(thread model.Thread) model.Thread
	GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread
	CheckParentPost(posts []model.Post, threadSlug string) bool
	ChangeThreadInDB(threadUpdate model.ThreadUpdate, oldThread model.Thread) model.Thread

	VoteForThreadInDB(slugOrId string, vote model.Vote) error
}

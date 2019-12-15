package repository

import (
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) IsForumInDB(slug string) bool {
	var count int
	err := r.db.Get(&count, `select count(*) from "`+forumTable+`" where slug=$1`, slug)
	checkErr(err)
	if count == 0 {
		return false
	}
	return true
}

func (r *DatabaseRepository) GetForumInDB(slug string) model.Forum {
	forumBySlug := model.Forum{}
	err := r.db.Get(&forumBySlug, `select * from "`+forumTable+`" where slug=$1`, slug)
	checkErr(err)
	return forumBySlug
}

func (r *DatabaseRepository) CreateForumInDB(forum model.Forum) model.Forum {
	userId := r.GetUserIdByName(forum.User)
	_, err := r.db.Exec(`insert into "`+forumTable+`" (slug, title, user_id) values ($1, $2, $3)`,
		forum.Slug, forum.Title, userId)
	checkErr(err)
	forum.Posts = 0
	forum.Treads = 0
	return forum
}

func (r *DatabaseRepository) GetForumUsersInDB(slug string, limit int, since string, desc bool) []model.User {
	var users []model.User
	threads := r.GetThreadsForumInDB(slug, limit, since, desc)
	for _, thread := range threads {
		user := r.GetUserById(thread.UserId)
		users = append(users, user)
	}

	posts := r.GetPostsForumInDB(slug, limit, since, desc)
	for _, post := range posts {
		user := r.GetUserById(post.UserId)
		users = append(users, user)
	}
	return users
}

func (r *DatabaseRepository) GetForumIdBySlug(slug string) int32 {
	var id int32
	err := r.db.Get(&id, `select id from "`+forumTable+`" where slug=$1 limit 1`, slug)
	checkErr(err)
	return id
}

func (r *DatabaseRepository) getForumById(id int32) model.Forum {
	forumById := model.Forum{}
	err := r.db.Get(&forumById, `select * from "`+forumTable+`" where id=$1`, id)
	checkErr(err)
	return forumById
}

func (r *DatabaseRepository) incForumDetails(field string, id int32) {
	forum := r.getForumById(id)
	var count int64 = 0
	if field == "posts" {
		count = forum.Posts
	} else if field == "threads" {
		count = int64(forum.Treads)
	}
	_, err := r.db.Exec(
		`update "`+forumTable+`" set $1=$2 where id=$3`,
		field, count + 1, id,
	)
	checkErr(err)
}
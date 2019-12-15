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
	err := r.db.Get(&forumBySlug, `select * from "`+forumTable+`" where slug=$1 limit 1`, slug)
	checkErr(err)
	return forumBySlug
}

func (r *DatabaseRepository) CreateForumInDB(forum model.Forum) model.Forum {
	forum.Posts = 0
	forum.Treads = 0
	userId := r.GetUserIdByName(forum.User)
	_, err := r.db.Exec(`insert into "`+forumTable+`" (slug, posts, title, threads, user_id) values ($1, $2, $3, $4, $5)`,
		forum.Slug, forum.Posts, forum.Title, forum.Treads, userId)
	checkErr(err)
	return forum
}

func (r *DatabaseRepository) GetForumUsersInDB(slug string, limit int, since string, desc bool) []model.User {
	// TODO
	return []model.User{}
}

func (r *DatabaseRepository) GetForumIdBySlug(slug string) int32 {
	var id int32
	err := r.db.Get(&id, `select id from "`+forumTable+`" where slug=$1 limit 1`, slug)
	checkErr(err)
	return id
}

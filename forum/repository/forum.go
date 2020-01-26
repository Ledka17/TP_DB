package repository

import (
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) IsForumInDB(slug string) bool {
	var count int
	err := r.db.Get(&count, `select count(*) from "`+forumTable+`" where lower(slug)=lower($1)`, slug)
	checkErr(err)
	if count == 0 {
		return false
	}
	return true
}

func (r *DatabaseRepository) GetForumInDB(slug string) model.Forum {
	forumBySlug := model.Forum{}
	err := r.db.Get(&forumBySlug, `select * from "`+forumTable+`" where lower(slug)=lower($1)`, slug)
	checkErr(err)
	return forumBySlug
}

func (r *DatabaseRepository) CreateForumInDB(forum model.Forum) model.Forum {
	_, err := r.db.Exec(`insert into "`+forumTable+`" (slug, title, author) values ($1, $2, $3)`,
		forum.Slug, forum.Title, forum.User)
	checkErr(err)
	return forum
}

func (r *DatabaseRepository) GetForumUsersInDB(slug string, limit int, since string, desc bool) []model.User {
	users := make([]model.User, 0)

	order := getOrder(desc)
	filterLimit := getFilterLimit(limit)
	filterSince := getFilterSinceByUserName(order, since)

	err := r.db.Select(&users, `select u.* from forum_user fu inner join "`+userTable+
		`" u on fu.user_nickname = u.nickname where lower(fu.forum)=lower($1) `+filterSince+ ` order by lower(u.nickname) `+
		order+filterLimit, slug)

	checkErr(err)
	return users
}

func (r *DatabaseRepository) GetForumIdBySlug(slug string) int32 {
	var id int32
	err := r.db.Get(&id, `select id from "`+forumTable+`" where lower(slug)=lower($1) limit 1`, slug)
	checkErr(err)
	return id
}

//func (r *DatabaseRepository) GetForumById(id int32) model.Forum {
//	forumById := model.Forum{}
//	err := r.db.Get(&forumById, `select * from "`+forumTable+`" where id=$1`, id)
//	checkErr(err)
//	return forumById
//}

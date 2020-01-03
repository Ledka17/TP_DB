package repository

import (
	"github.com/Ledka17/TP_DB/model"
	"strconv"
)

func (r *DatabaseRepository) IsThreadInDB(slugOrId string) bool {
	count := 0
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		id = -1
	}
	if slugOrId != "" {
		err := r.db.Get(&count, `select count(*) from "`+threadTable+`" where lower(slug)=lower($1) or id=$2`, slugOrId, id)
		checkErr(err)
	}
	if count > 0 {
		return true
	}
	return false
}

func (r *DatabaseRepository) GetThreadInDB(slugOrId string) model.Thread {
	var thread model.Thread
	id, _ := strconv.Atoi(slugOrId)
	err := r.db.Get(&thread, `select * from "`+threadTable+`" where lower(slug)=lower($1) or id=$2 limit 1`, slugOrId, id)
	checkErr(err)
	thread.Forum = r.GetForumById(thread.ForumId).Slug
	thread.Author = r.GetUserById(thread.UserId).Nickname
	return thread
}

func (r *DatabaseRepository) CreateThreadInDB(forumSlug string, thread model.Thread) model.Thread {
	thread.Forum = r.GetForumInDB(forumSlug).Slug
	//thread.Created = time.Now()
	thread.ForumId = r.GetForumIdBySlug(thread.Forum)
	thread.UserId = r.GetUserInDB(thread.Author).Id

	err := r.db.QueryRow(`insert into "`+threadTable+
		`" (title, slug, user_id, message, created, forum_id) values ($1, $2, $3, $4, $5, $6) returning id`,
		thread.Title, thread.Slug, thread.UserId, thread.Message, thread.Created, thread.ForumId).Scan(&thread.Id)
	checkErr(err)
	r.incForumDetails("threads", thread.ForumId)
	return thread
}

func (r *DatabaseRepository) GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread {
	// TODO where incorrect syntax
	threads := make([]model.Thread, 0, limit)
	forumId := r.GetForumIdBySlug(forumSlug)
	order := getOrder(desc)
	filterLimit := getFilterLimit(limit)
	filterSince := getFilterSince(order, since)

	err := r.db.Select(&threads, `select * from "`+threadTable+`" where forum_id=$1 `+filterSince+ ` order by $2 `+filterLimit,
		forumId, order,
	)
	checkErr(err)
	return threads
}

func (r *DatabaseRepository) CheckParentPost(posts []model.Post) bool {
	var parentsForCheck, children []int64
	for _, post := range posts { // выгружаем всех родителей и детей
		if !have(post.Parent, parentsForCheck) {
			parentsForCheck = append(parentsForCheck, post.Parent)
		}
		if !have(post.Id, children) {
			children = append(children, post.Id)
		}
	}
	for _, parent := range parentsForCheck { // проверяем есть ли родитель
		if parent != 0 && !have(parent, children) && !r.IsPostInDB(int(parent)) {
			return false
		}
	}
	return true
}

func (r *DatabaseRepository) ChangeThreadInDB(threadUpdate model.ThreadUpdate, slugOrId string) model.Thread {
	thread := r.GetThreadInDB(slugOrId)
	if threadUpdate.Title != "" {
		thread.Title = threadUpdate.Title
		_, err := r.db.Exec(
			`update "`+threadTable+`" set title=$1 where id=$2`,
			threadUpdate.Title, thread.Id,
		)
		checkErr(err)
	}
	if threadUpdate.Message != "" {
		thread.Message = threadUpdate.Message
		_, err := r.db.Exec(
			`update "`+threadTable+`" set message=$1 where id=$2`,
			threadUpdate.Message, thread.Id,
		)
		checkErr(err)
	}
	return thread
}

func have(elem int64, array []int64) bool {
	for _, current := range array {
		if current == elem {
			return true
		}
	}
	return false
}
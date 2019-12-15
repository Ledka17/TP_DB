package repository

import (
	"github.com/Ledka17/TP_DB/model"
	"time"
)

func (r *DatabaseRepository) IsThreadInDB(slugOrId string) bool {
	var count int
	err := r.db.Get(&count, `select count(*) from "`+threadTable+`" where slug=$1 or id=$1`, slugOrId)
	checkErr(err)
	if count > 0 {
		return true
	}
	return false
}

func (r *DatabaseRepository) GetThreadInDB(slugOrId string) model.Thread {
	var thread, emptyThread model.Thread
	err := r.db.Get(&thread, `select * from "`+threadTable+`" where slug=$1`, slugOrId)
	checkErr(err)
	if thread != emptyThread {
		return thread
	}

	err = r.db.Get(&thread, `select count(*) from "`+threadTable+`" where id=$1`, slugOrId)
	checkErr(err)
	return thread
}

func (r *DatabaseRepository) CreateThreadInDB(forumSlug string, thread model.Thread) model.Thread {
	thread.Forum = forumSlug
	thread.Created = time.Now()
	thread.ForumId = r.GetForumIdBySlug(forumSlug)

	_, err := r.db.Exec(`insert into "`+threadTable+
		`" (title, slug, author, message, created, forum_id) values ($1, $2, $3, $4, $5, $6)`,
		thread.Title, thread.Slug, thread.Author, thread.Message, thread.Created, thread.ForumId)
	checkErr(err)
	r.incForumDetails("threads", thread.ForumId)
	return thread
}

func (r *DatabaseRepository) GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread {
	threads := make([]model.Thread, 0, limit)
	forumId := r.GetForumIdBySlug(forumSlug)
	order := getOrder(desc)
	filterLimit := getFilterLimit(limit)
	filterSince := getFilterSince(order, since)

	err := r.db.Select(&threads, `select * from "`+threadTable+`" where forum_id=$1 `+filterSince+ `order by $2 `+filterLimit,
		forumId, since, order,
	)
	checkErr(err)
	return threads
}

func (r *DatabaseRepository) CheckParentPost(posts []model.Post) bool {
	// TODO
	var parentsForCheck, children []int64
	for _, post := range posts { // выгружаем всех родителей и детей
		parentsForCheck = append(parentsForCheck, post.Parent)
		children = append(children, post.Id)
	}
	//for _, child := range children { // проверяем есть у ребенка родитель
	//}
	return false
}

func (r *DatabaseRepository) ChangeThreadInDB(threadUpdate model.ThreadUpdate, slugOrId string) model.Thread {
	thread := r.GetThreadInDB(slugOrId)
	thread.Title = threadUpdate.Title
	thread.Message = threadUpdate.Message
	_, err := r.db.Exec(
		`update "`+threadTable+`" set message=$1, title=$2 where id=$3`,
		threadUpdate.Message, threadUpdate.Title, thread.Id,
	)
	checkErr(err)
	return thread
}

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
	var thread, emptyThread model.Thread
	err := r.db.Get(&thread, `select * from "`+threadTable+`" where lower(slug)=lower($1) limit 1`, slugOrId)
	checkErr(err)
	//fmt.Println("select")
	if thread == emptyThread {
		id, _ := strconv.Atoi(slugOrId)
		thread = r.GetThreadById(id)
	}
	return thread
}

func (r *DatabaseRepository) GetThreadById(id int) model.Thread {
	var thread model.Thread
	err := r.db.Get(&thread, `select * from "`+threadTable+`" where id=$1 limit 1`, id)
	checkErr(err)
	return thread
}

func (r *DatabaseRepository) CreateThreadInDB(thread model.Thread) model.Thread {
	tx, err := r.db.Beginx()
	defer tx.Rollback()
	err = tx.QueryRow(`insert into "`+threadTable+
		`" (title, slug, user_id, message, created, forum_id, author, forum) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`,
		thread.Title, thread.Slug, thread.UserId, thread.Message, thread.Created, thread.ForumId, thread.Author, thread.Forum).
		Scan(&thread.Id)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	checkErr(err)
	//r.incForumDetails("threads", thread.ForumId)
	return thread
}

func (r *DatabaseRepository) GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread {
	threads := make([]model.Thread, 0)
	order := getOrder(desc)
	filterLimit := getFilterLimit(limit)
	filterSince := getFilterSince(order, since)

	err := r.db.Select(&threads, `select distinct * from "`+threadTable+`" where forum=$1 `+filterSince+
		` order by created `+order+filterLimit,
		forumSlug,
	)
	checkErr(err)
	return threads
}

func (r *DatabaseRepository) CheckParentPost(posts []model.Post, threadSlug string) bool {
	var parentsForCheck []model.Post
	var children []int64
	emptyPost := model.Post{}
	threadId := r.GetThreadInDB(threadSlug).Id
	for _, post := range posts { // выгружаем всех родителей и детей
		parentsForCheck = append(parentsForCheck, post)
		children = append(children, post.Id)
	}
	for _, post := range parentsForCheck { // проверяем есть ли родитель
		if post.Parent != 0 && !have(post.Parent, children) {
			parentPost := r.getPostById(int(post.Parent))
			if parentPost == emptyPost || parentPost.ThreadId != threadId {
				return false
			}
		}
	}
	return true
}

func (r *DatabaseRepository) ChangeThreadInDB(threadUpdate model.ThreadUpdate, thread model.Thread) model.Thread {
	tx, err := r.db.Beginx()
	checkErr(err)
	defer tx.Rollback()

	if threadUpdate.Title != "" || threadUpdate.Message != "" {
		thread.Title = threadUpdate.Title
		err := tx.QueryRow(
			`update "`+threadTable+`" set title = coalesce(nullif($1, ''), title), message = coalesce(nullif($2, ''), message) where id=$3 returning title, message`,
			threadUpdate.Title, threadUpdate.Message, thread.Id,
		).Scan(&thread.Title, &thread.Message)
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}
	return thread
}

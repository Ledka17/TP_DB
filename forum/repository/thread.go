package repository

import (
	"github.com/Ledka17/TP_DB/model"
	"time"
)

func (r *DatabaseRepository) IsThreadInDB(slugOrId string) bool {
	var countSlug int
	err := r.db.Get(&countSlug, `select count(*) from "`+threadTable+`" where slug=$1`, slugOrId)
	checkErr(err)
	if countSlug != 0 {
		return true
	}

	var countId int
	err = r.db.Get(&countId, `select count(*) from "`+threadTable+`" where id=$1`, slugOrId)
	checkErr(err)
	if countId != 0 {
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
	thread.Votes = 0
	thread.Id = r.GetNextThreadId()
	forumId := r.GetForumIdBySlug(forumSlug)

	_, err := r.db.Exec(`insert into "`+threadTable+
		`" (title, slug, author, message, votes, created, forum_id, id) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		thread.Title, thread.Slug, thread.Author, thread.Message, thread.Votes, thread.Created, forumId, thread.Id)
	checkErr(err)
	return model.Thread{}
}

func (r *DatabaseRepository) GetThreadsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Thread {
	// TODO
	return []model.Thread{}
}

func (r *DatabaseRepository) CheckParentPost(posts []model.Post) bool {
	// TODO
	var parents, children []int64
	for _, parent := range posts { // выгружаем всех родителей и детей
		parents = append(parents, parent.Parent)
		children = append(children, parent.Id)
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

func (r *DatabaseRepository) GetNextThreadId() int32 {
	var nextId int32 = -1
	err := r.db.Get(&nextId, `select max(id) from "`+threadTable)
	checkErr(err)
	return nextId + 1
}
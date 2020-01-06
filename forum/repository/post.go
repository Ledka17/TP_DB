package repository

import (
	"github.com/Ledka17/TP_DB/model"
	"time"
)

func (r *DatabaseRepository) IsPostInDB(id int) bool {
	var count int
	err := r.db.Get(&count, `select count(*) from "`+postTable+`" where id=$1`, id)
	checkErr(err)
	if count > 0 {
		return true
	}
	return false
}

func (r *DatabaseRepository) GetPostInDB(id int) model.Post {
	return r.getPostById(id)
}

func (r *DatabaseRepository) GetPostsInDB(threadSlugOrId string, limit int, since int, sort string, desc bool) []model.Post {
	order := getOrder(desc)
	threadId := r.GetThreadInDB(threadSlugOrId).Id
	switch sort {
	case "flat", "":
		return r.getPostsFlat(threadId, limit, since, order)
	case "tree":
		return r.getPostsTree(threadId, limit, since, order)
	case "parent_tree":
		return r.getPostsParentTree(threadId, limit, since, order)
	}
	return []model.Post{}
}

func (r *DatabaseRepository) ChangePostInDB(id int, update model.PostUpdate) model.Post {
	post := r.getPostById(id)
	if update.Message != "" && update.Message != post.Message{
		_, err := r.db.Exec(
			`update "`+postTable+`" set message=$1, isEdited=True where id=$2`,
			update.Message, id,
		)
		checkErr(err)
		post.IsEdited = true
		post.Message = update.Message
	}
	return post
}

func (r *DatabaseRepository) CreatePostsInDB(posts []model.Post, threadSlugOrId string) []model.Post {
	created := time.Now().Format(time.RFC3339)
	for i, post := range posts {
		curThread := r.GetThreadInDB(threadSlugOrId)
		post.ThreadId = curThread.Id
		post.ForumId = curThread.ForumId
		post.UserId = r.GetUserIdByName(post.Author)
		post.Created = created
		post.Forum = curThread.Forum

		err := r.db.QueryRow(`insert into "`+postTable+`" (parent, message, created, user_id, forum_id, thread_id, forum, author) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`,
			post.Parent, post.Message, post.Created, post.UserId, post.ForumId, post.ThreadId, post.Forum, post.Author).Scan(&post.Id)
		checkErr(err)
		r.incForumDetails("posts", post.ForumId)

		posts[i] = post
	}
	return posts
}

func (r *DatabaseRepository) getPostById(id int) model.Post {
	var post model.Post
	err := r.db.Get(&post, `select * from "`+postTable+`" where id=$1`, id)
	checkErr(err)
	return post
}

func (r *DatabaseRepository) getPostsFlat(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	filterId := getFilterId(order, since)
	filterLimit := getFilterLimit(limit)

	err := r.db.Select(&posts, `select * from "`+postTable+`" where thread_id=$1 `+filterId+` order by id `+order+filterLimit,
		threadId,
	)
	checkErr(err)
	return posts
}

func (r *DatabaseRepository) getPostsTree(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	filterLimit := getFilterLimit(limit)

	recursiveQuery := `with recursive r as (select *, CAST (id AS VARCHAR (50)) as PATH, CAST (id AS VARCHAR (50)) as OLDPATH, id as LEVEL0 from "`+
		postTable+`" where thread_id=$1 and parent=0`+
		`union select p.*, CAST ( r.PATH ||'->'|| p.id+10000000 AS VARCHAR(50)) as PATH, CAST (r.PATH  AS VARCHAR (50)) as OLDPATH, LEVEL0 as LEVEL0 from "`+
		postTable+ `" p inner join r on p.parent = r.id) `

	if since == -1 {
		err := r.db.Select(&posts, recursiveQuery+
			`select id, thread_id, parent, author, created, forum, isedited, message from r order by LEVEL0 `+
			order+`, path `+order+filterLimit,
			threadId,
		)
		checkErr(err)
	} else {
		sign := ">"
		resTable := `(select row_number() over(order by level0 `+order+`, path `+order+
			`) as n, * from r order by LEVEL0 `+order+`, path `+order+`) r `
		err := r.db.Select(&posts, recursiveQuery+
			`select r.id, r.thread_id, r.parent, r.author, r.created, r.forum, r.isedited, r.message from `+resTable+
			` where n `+sign+` (select r.n from `+resTable+` where id=$2)`+filterLimit,
			threadId, since,
		)
		checkErr(err)
	}

	return posts
}

func (r *DatabaseRepository) getPostsParentTree(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	filterId := ""
	filterLimit := getFilterLimit(limit)

	if since != -1 {
		parentSince := r.getPostById(since).Parent
		filterId = getFilterId(order, int(parentSince))
	}

	recursiveQuery := `with recursive r as ( (select *, CAST (id AS VARCHAR (50)) as PATH, id as LEVEL0 from"`+
		postTable+`" where thread_id=$1 and parent=0 `+filterId+`order by id `+order+filterLimit+
		`) union (select p.*, CAST ( r.PATH ||'->'|| p.id+10000000 AS VARCHAR(50)) as PATH, LEVEL0 as LEVEL0 from "`+
		postTable+ `" p inner join r on p.parent = r.id) ) `
	err := r.db.Select(&posts, recursiveQuery+
		`select id, thread_id, parent, author, created, forum, isedited, message from r order by LEVEL0 `+
		order+`, path`,
		threadId,
	)
	checkErr(err)
	return posts
}

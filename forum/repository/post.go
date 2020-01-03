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
	posts := make([]model.Post, 0, limit)
	order := getOrder(desc)
	threadId := r.GetThreadInDB(threadSlugOrId).Id
	switch sort {
	case "flat", "":
		posts = r.getPostsFlat(threadId, limit, since, order)
	case "tree":
		// TODO
		return []model.Post{}
	case "parent_tree":
		// TODO
		return []model.Post{}
	}
	return posts
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
	//post.Author = r.GetUserById(post.UserId).Nickname
	//post.Forum = r.GetForumById(post.ForumId).Slug
	return post
}

func (r *DatabaseRepository) getPostsFlat(threadId int32, limit int, since int, order string) []model.Post {
	var posts []model.Post

	filterId := getFilterId(order, since)
	filterLimit := getFilterLimit(limit)

	err := r.db.Select(&posts, `select * from "`+postTable+`" where 1=1`+filterId+` order by $1 `+filterLimit,
		order,
	)
	checkErr(err)
	return posts
}

func (r *DatabaseRepository) GetPostsForumInDB(forumSlug string, limit int, since string, desc bool) []model.Post {
	posts := make([]model.Post, 0, limit)
	forumId := r.GetForumIdBySlug(forumSlug)
	order := getOrder(desc)
	filterLimit := getFilterLimit(limit)
	filterSince := getFilterSince(order, since)

	err := r.db.Select(&posts, `select * from "`+postTable+`" where forum_id=$1 `+filterSince+ `order by $2 `+filterLimit,
		forumId, order,
	)
	checkErr(err)
	return posts
}
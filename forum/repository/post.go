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

func (r *DatabaseRepository) GetPostInDB(id int, related []string) model.PostFull {
	var postFull model.PostFull
	post := r.getPostById(id)
	postFull.Post = post
	if checkInRelated("user", related) {
		postFull.Author = r.GetUserInDB(post.Author)
	}
	if checkInRelated("forum", related) {
		postFull.Forum = r.getForumById(post.ForumId)
	}
	if checkInRelated("thread", related) {
		postFull.Thread = r.GetThreadInDB(string(post.ThreadId))
	}
	return postFull
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
	post.Message = update.Message
	_, err := r.db.Exec(
		`update "`+postTable+`" set message=$1 isEdited=True where id=$2`,
		post.Message, post.Id,
	)
	checkErr(err)
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

		err := r.db.QueryRow(`insert into "`+postTable+`" (parent, message, created, user_id, forum_id, thread_id) values ($1, $2, $3, $4, $5, $6) returning id`,
			post.Parent, post.Message, post.Created, post.UserId, post.ForumId, post.ThreadId).Scan(&post.Id)
		checkErr(err)
		r.incForumDetails("posts", post.ForumId)

		post.Forum = curThread.Forum
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
	var posts []model.Post

	filterId := getFilterId(order, since)
	filterLimit := getFilterLimit(limit)

	err := r.db.Select(&posts, `select * from "`+postTable+`"`+filterId+`order by $1 `+filterLimit,
		order,
	)
	checkErr(err)
	return posts
}

func checkInRelated(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
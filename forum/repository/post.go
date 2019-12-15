package repository

import "github.com/Ledka17/TP_DB/model"

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
	// TODO
	return model.PostFull{}
}

func (r *DatabaseRepository) GetPostsInDB(threadSlugOrId string, limit int, since int, sort string, desc bool) []model.Post {
	// TODO
	return []model.Post{}
}

func (r *DatabaseRepository) ChangePostInDB(id int, update model.PostUpdate) model.Post {
	// TODO
	return model.Post{}
}

func (r *DatabaseRepository) CreatePostsInDB(posts []model.Post) []model.Post {
	for _, post := range posts {
		post.IsEdited = false
		forumId := r.GetForumIdBySlug(post.Forum)
		_, err := r.db.Exec(`insert into "`+postTable+`" (id, parent, message, created, author, forum_id, thread, isEdited) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
			post.Id, post.Parent, post.Message, post.Created, post.Author, forumId, post.Tread, post.IsEdited)
		checkErr(err)
	}
	return posts
}

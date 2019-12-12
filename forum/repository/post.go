package repository

import "TP_DB/model"

func (r *DatabaseRepository) IsPostInDB(id int) bool {
	// TODO
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
	// TODO
	return []model.Post{}
}

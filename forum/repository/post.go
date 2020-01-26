package repository

import (
	"database/sql"
	"fmt"
	"github.com/Ledka17/TP_DB/model"
	"strings"
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

func (r *DatabaseRepository) GetPostsInDB(threadId int32, limit int, since int, sort string, desc bool) []model.Post {
	order := getOrder(desc)
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

func (r *DatabaseRepository) ChangePostInDB(id int, update model.PostUpdate) (model.Post, error) {
	post := r.getPostById(id)
	if update.Message != "" && update.Message != post.Message {
		_, err := r.db.Exec(
			`update "`+postTable+`" set message=$1, isEdited=True where id=$2`,
			update.Message, id,
		)
		if err == sql.ErrNoRows || err != nil {
			return model.Post{}, err
		}
		post.IsEdited = true
		post.Message = update.Message
	}
	return post, nil
	////post := r.getPostById(id)
	//tx, _ := r.db.Beginx()
	//defer tx.Rollback()
	//
	//if update.Message != "" {
	//	post := model.Post{}
	//	err := r.db.QueryRow(
	//		`update "`+postTable+`" set message=$1, isEdited= case when message != $1 then True else False end where id=$2`+
	//			`returning id, author, created, forum, isedited, message, parent, thread_id`,
	//		update.Message, id,
	//	).Scan(&post.Id, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.ThreadId)
	//	fmt.Println(err)
	//	if err == sql.ErrNoRows || err != nil {
	//		tx.Rollback()
	//		return model.Post{}, err
	//	}
	//	tx.Commit()
	//	return post, nil
	//	//post.IsEdited = true
	//	//post.Message = update.Message
	//}
	//post := r.getPostById(id)
	//tx.Commit()
	//return post, nil
}

func (r *DatabaseRepository) CreatePostsInDB(posts []model.Post, threadSlugOrId string) ([]model.Post, error) {
	created := time.Now().Format(time.RFC3339)
	tx, err := r.db.Beginx()
	checkErr(err)
	defer tx.Rollback()

	curThread := r.GetThreadInDB(threadSlugOrId)
	emptyThread := model.Thread{}
	if curThread == emptyThread {
		tx.Rollback()
		return []model.Post{}, fmt.Errorf("thread not found")
	}

	if len(posts) == 0 {
		tx.Commit()
		return posts, nil
	}

	columns := 6
	queryValues := make([]string, len(posts))
	args := make([]interface{}, 0, len(posts)*columns)

	for i, post := range posts {
		if post.Author == "" || !r.IsUserInDB(post.Author, "") {
			tx.Rollback()
			return []model.Post{}, fmt.Errorf("user not found")
		}
		post.ThreadId = curThread.Id
		post.Created = created
		post.Forum = curThread.Forum

		if post.Parent > 0 {
			parentPost := r.GetPostInDB(int(post.Parent))
			emptyPost := model.Post{}
			if parentPost == emptyPost || parentPost.ThreadId != post.ThreadId {
				tx.Rollback()
				return posts, fmt.Errorf("conflicts in posts")
			}
		}

		values := fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d)`, i*columns+1,i*columns+2,i*columns+3,i*columns+4,i*columns+5,i*columns+6)
		args = append(args, post.Parent, post.Message, post.Created, post.ThreadId, post.Forum, post.Author)

		queryValues[i] = values
		//err := tx.QueryRow(`insert into "`+postTable+`" (parent, message, created, thread_id, forum, author) values ($1, $2, $3, $4, $5, $6) returning id`,
		//	post.Parent, post.Message, post.Created, post.ThreadId, post.Forum, post.Author).Scan(&post.Id)
		//if err != nil {
		//	tx.Rollback()
		//	return posts, fmt.Errorf("conflicts in posts")
		//}

		posts[i] = post
	}
	queryStart := `INSERT INTO "`+postTable+`" (parent, message, created, thread_id, forum, author) values `
	queryEnd := ` returning id`

	rows, err := tx.Query(queryStart + strings.Join(queryValues, ", ") + queryEnd, args...)

	//fmt.Println(err)
	//fmt.Println(rows)
	if err != nil {
		tx.Rollback()
		return posts, fmt.Errorf("conflicts in posts")
	}
	for i := range posts {
		if rows.Next() {
			err1 := rows.Scan(&posts[i].Id)
			fmt.Println("test", posts[i])
			fmt.Println("err", err1)
		}
	}
	rows.Close()
	//for i, _ := range posts {
	//	posts[i].Id = postsId[i]
	//}

	tx.Commit()
	return posts, nil
}

func (r *DatabaseRepository) getPostById(id int) model.Post {
	var post model.Post
	err := r.db.Get(&post, `select * from "`+postTable+`" where id=$1`, id)
	checkErr(err)
	return post
}

func (r *DatabaseRepository) getPostsFlat(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	tx, err := r.db.Beginx()
	defer tx.Rollback()
	filterId := getFilterId(order, since)
	filterLimit := getFilterLimit(limit)

	err = tx.Select(&posts, `select * from "`+postTable+`" where thread_id=$1 `+filterId+` order by id `+order+filterLimit,
		threadId,
	)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return posts
}

func (r *DatabaseRepository) getPostsTree(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	tx, err := r.db.Beginx()
	defer tx.Rollback()

	filterLimit := getFilterLimit(limit)
	filterSince := ""

	if since != -1 {
		sign := getSign(order)
		filterSince = fmt.Sprintf(" and path%s(select path from post where id=%d) ", sign, since)
	}

	err = tx.Select(&posts,`select * from post where thread_id=$1`+filterSince+
		`order by path `+order+filterLimit,
		threadId,
	)

	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return posts
}

func (r *DatabaseRepository) getPostsParentTree(threadId int32, limit int, since int, order string) []model.Post {
	posts := make([]model.Post, 0)

	tx, err := r.db.Beginx()
	defer tx.Rollback()

	filterId := ""
	filterLimit := getFilterLimit(limit)

	if since != -1 {
		parentSince := r.getPostById(since).Parent
		filterId = getFilterId(order, int(parentSince))
	}

	err = tx.Select(&posts, `select * from post where left(path,8) in (select distinct parents.path from post as parents where `+
		` parents.thread_id=$1 and parents.parent=0`+filterId+` order by path `+order+
		filterLimit+`) order by left(path,8) `+order+`, path;`,
		threadId)

	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return posts
}

package repository

import (
	"github.com/Ledka17/TP_DB/forum"
	"github.com/jmoiron/sqlx"
)

const (
	userTable = "user"
	forumTable = "forum"
	threadTable = "thread"
	postTable = "post"
	voteTable = "vote"
)

type DatabaseRepository struct {
	db                 *sqlx.DB
}

func NewDatabaseRepository(db *sqlx.DB) forum.Repository {
	return &DatabaseRepository{
		db,
	}
}

package repository

import (
	"TP_DB/forum"
	"github.com/jmoiron/sqlx"
	"log"
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

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
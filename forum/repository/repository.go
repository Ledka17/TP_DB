package repository

import (
	"database/sql"
	"fmt"
	"github.com/Ledka17/TP_DB/forum"
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
	if err != sql.ErrNoRows && err != nil {
		log.Panic(err)
	}
}

func getOrder(desc bool) string {
	if desc {
		return "desc"
	}
	return "asc"
}

func getFilterLimit (limit int) string {
	filterLimit := ""
	if limit > 0 {
		filterLimit = fmt.Sprintf(" limit %d ", limit)
	}
	return filterLimit
}

func getFilterId (order string, id int) string {
	filterId := fmt.Sprintf(" id > %d ", id)
	if order == "desc" {
		filterId = fmt.Sprintf(" id < %d ", id)
	}
	return filterId
}

func getFilterSince(order string, since string) string{
	filterSince := fmt.Sprintf(" created > '%s' ", since)
	if order == "desc" {
		filterSince = fmt.Sprintf(" created < '%s' ", since)
	}
	return filterSince
}
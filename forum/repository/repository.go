package repository

import (
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
	if err != nil {
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
		filterLimit = fmt.Sprintf("limit %d", limit)
	}
	return filterLimit
}

func getFilterId (order string, id int) string {
	filterId := fmt.Sprintf(" where id > %d ", id)
	if order == "desc" {
		filterId = fmt.Sprintf(" where id < %d ", id)
	}
	return filterId
}

func getFilterSince(order string, since string) string{
	filterSince := fmt.Sprintf(" where created > %s ", since)
	if order == "desc" {
		filterSince = fmt.Sprintf(" where id < %s ", since)
	}
	return filterSince
}
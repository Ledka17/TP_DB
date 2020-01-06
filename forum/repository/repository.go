package repository

import (
	"database/sql"
	"fmt"
	"github.com/Ledka17/TP_DB/forum"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
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
	if id == -1 {
		return ""
	}
	filterId := fmt.Sprintf(" and id > %d ", id)
	if order == "desc" {
		filterId = fmt.Sprintf(" and id < %d ", id)
	}
	return filterId
}

func getFilterSince(order string, since string) string{
	if since == "" {
		return ""
	}

	log.Println("time =", since)
	x,_ := time.Parse("2019-09-13T16:51:23.112Z", since)
	log.Println("parse =", x)

	filterSince := fmt.Sprintf(" and created >= '%s' ", since)
	if order == "desc" {
		filterSince = fmt.Sprintf(" and created <= '%s' ", since)
	}
	return filterSince
}

func getFilterSinceByUserName(order string, since string) string{
	if since == "" {
		return ""
	}
	filterSince := fmt.Sprintf(" and lower(u.nickname) > lower('%s') ", since)
	if order == "desc" {
		filterSince = fmt.Sprintf(" and lower(u.nickname) < lower('%s') ", since)
	}
	return filterSince
}
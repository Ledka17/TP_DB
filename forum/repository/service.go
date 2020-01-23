package repository

import (
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) CleanUp()  {
	_, err := r.db.Exec(`truncate "`+forumTable+`", "`+postTable+`", "`+threadTable+`", "`+voteTable+`", "`+userTable+`"`)
	checkErr(err)
}

func (r *DatabaseRepository) GetStatusDB() model.Status {
	var status model.Status
	tx, err := r.db.Beginx()
	checkErr(err)

	defer tx.Rollback()

	var userCount int32
	err = tx.QueryRow(`select count(*) from "`+userTable+`"`).Scan(&userCount)
	checkErr(err)

	var threadCount int32
	err = tx.QueryRow(`select count(*) from "`+threadTable+`"`).Scan(&threadCount)
	checkErr(err)

	var postCount int64
	err = tx.QueryRow(`select count(*) from "`+postTable+`"`).Scan(&postCount)
	checkErr(err)

	var forumCount int32
	err = tx.QueryRow(`select count(*) from "`+forumTable+`"`).Scan(&forumCount)
	checkErr(err)

	err = tx.Commit()
	checkErr(err)

	status.User = userCount
	status.Thread = threadCount
	status.Post = postCount
	status.Forum = forumCount

	return status
}

func (r *DatabaseRepository) countRecords(tableName string) int64 {
	var countOfRecords int64 = 0
	err := r.db.Get(&countOfRecords, `select count(*) from "`+tableName+`"`)
	checkErr(err)
	return countOfRecords
}
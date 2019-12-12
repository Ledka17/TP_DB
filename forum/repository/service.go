package repository

import "github.com/Ledka17/TP_DB/model"

func (r *DatabaseRepository) CleanUp()  {
	r.deleteTable(forumTable)
	r.deleteTable(voteTable)
	r.deleteTable(threadTable)
	r.deleteTable(userTable)
	r.deleteTable(postTable)
}

func (r *DatabaseRepository) GetStatusDB() model.Status {
	var status model.Status
	status.User = int32(r.countRecords(userTable))
	status.Thread = int32(r.countRecords(threadTable))
	status.Post = r.countRecords(postTable)
	status.Forum = int32(r.countRecords(forumTable))

	return status
}

func (r *DatabaseRepository) deleteTable(tableName string) {
	_, err := r.db.Exec(`delete from "`+tableName)
	checkErr(err)
}

func (r *DatabaseRepository) countRecords(tableName string) int64 {
	var countOfRecords int64 = 0
	err := r.db.Get(&countOfRecords, `select count(*) from "`+tableName)
	checkErr(err)
	return countOfRecords
}
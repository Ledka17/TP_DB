package repository

import "github.com/Ledka17/TP_DB/model"

func (r *DatabaseRepository) VoteForThreadInDB(threadSlugOrId string, vote model.Vote) model.Thread {
	thread := r.GetThreadInDB(threadSlugOrId)
	_, err := r.db.Exec(`insert into "`+voteTable+
		`" (thread_id, nickname, voice) values ($1, $2, $3)`,
		thread.Id, vote.Nickname, vote.Voice)
	checkErr(err)
	thread.Votes += vote.Voice
	_, err = r.db.Exec(
		`update "`+threadTable+`" set votes=$1 where id=$3`,
		thread.Votes,  thread.Id,
	)
	checkErr(err)
	return thread
}

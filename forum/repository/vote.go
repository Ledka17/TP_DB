package repository

import "github.com/Ledka17/TP_DB/model"

func (r *DatabaseRepository) VoteForThreadInDB(threadSlugOrId string, vote model.Vote) model.Thread {
	var emptyVote model.Vote
	thread := r.GetThreadInDB(threadSlugOrId)
	voteGet := r.getVoteInDB(thread.Id, vote.Nickname)

	if voteGet == emptyVote {
		_, err := r.db.Exec(`insert into "`+voteTable+
			`" (thread_id, nickname, voice) values ($1, $2, $3)`,
			thread.Id, vote.Nickname, vote.Voice)
		checkErr(err)
		thread.Votes += vote.Voice
	} else {
		_, err := r.db.Exec(`update "`+voteTable+
			`" set voice=$1 where id=$2`,
			vote.Voice, voteGet.Id,
		)
		checkErr(err)
		thread.Votes += vote.Voice - voteGet.Voice
	}

	_, err := r.db.Exec(
		`update "`+threadTable+`" set votes=$1 where id=$2`,
		thread.Votes, thread.Id,
	)
	checkErr(err)
	return thread
}

func (r *DatabaseRepository) getVoteInDB(threadId int32, nickname string) model.Vote {
	var vote model.Vote
	err := r.db.Get(&vote, `select * from "`+voteTable+`" where thread_id=$1 and nickname=$2 limit 1`,
		threadId, nickname)
	checkErr(err)
	return vote
}

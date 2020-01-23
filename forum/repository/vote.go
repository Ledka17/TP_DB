package repository

import "github.com/Ledka17/TP_DB/model"

func (r *DatabaseRepository) VoteForThreadInDB(thread model.Thread, vote model.Vote) model.Thread {
	var emptyVote model.Vote
	foundVote := r.getVoteInDB(thread.Id, vote.Nickname)

	tx, err := r.db.Beginx()
	defer tx.Rollback()

	if foundVote == emptyVote {
		_, err := tx.Exec(`insert into "`+voteTable+
			`" (thread_id, nickname, voice) values ($1, $2, $3)`,
			thread.Id, vote.Nickname, vote.Voice)
		if err != nil {
			tx.Rollback()
		}
		thread.Votes += vote.Voice
	} else {
		_, err := tx.Exec(`update "`+voteTable+
			`" set voice=$1 where id=$2`,
			vote.Voice, foundVote.Id,
		)
		if err != nil {
			tx.Rollback()
		}
		thread.Votes += vote.Voice - foundVote.Voice
	}

	//_, err = tx.Exec(
	//	`update "`+threadTable+`" set votes=$1 where id=$2`,
	//	thread.Votes, thread.Id,
	//)
	//if err != nil {
	//	tx.Rollback()
	//}

	err = tx.Commit()
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

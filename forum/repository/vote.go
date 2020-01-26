package repository

import (
	"fmt"
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) VoteForThreadInDB(slugOrId string, vote model.Vote) error {
	thread := r.GetThreadInDB(slugOrId)
	emptyThread := model.Thread{}

	if thread == emptyThread {
		return fmt.Errorf("thread not found")
	}

	userIs := r.IsUserInDB(vote.Nickname, "")
	if !userIs {
		return fmt.Errorf("user not found")
	}
	_, err := r.db.Exec(`insert into "`+voteTable+`" (thread_id, nickname, voice) VALUES ($1, $2, $3)
		ON CONFLICT ON CONSTRAINT vote_thread_id_nickname_key DO
		UPDATE SET voice=$3 WHERE vote.thread_id=$1 AND lower(vote.nickname)=lower($2)`,
		thread.Id, vote.Nickname, vote.Voice)
	checkErr(err)
	//foundVote := r.getVoteInDB(thread.Id, user.Nickname)
	//emptyVote := model.Vote{}
	//
	//if foundVote == emptyVote {
	//	_, err := r.db.Exec(`insert into "`+voteTable+`" (thread_id, nickname, voice) VALUES ($1, $2, $3)`,
	//		thread.Id, vote.Nickname, vote.Voice)
	//	checkErr(err)
	//} else {
	//	_, err := r.db.Exec(`update "`+voteTable+`" set voice=$3 where thread_id=$1 and nickname=$2`,
	//		thread.Id, vote.Nickname, vote.Voice)
	//	checkErr(err)
	//}

	return nil
}

func (r *DatabaseRepository) getVoteInDB(threadId int32, nickname string) model.Vote {
	var vote model.Vote
	err := r.db.Get(&vote, `select * from "`+voteTable+`" where thread_id=$1 and nickname=$2 limit 1`,
		threadId, nickname)
	checkErr(err)
	return vote
}

package repository

import (
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) IsUserInDB(nickname string, email string) bool {
	var count int
	err := r.db.Get(&count, `select count(*) from "`+userTable+`" where nickname=$1 or email=$2`, nickname, email)
	checkErr(err)
	if count != 0 {
		return true
	}
	return false
}

func (r *DatabaseRepository) GetUserInDB(nickname string, email string) model.User {
	var user model.User
	err := r.db.Get(&user, `select * from "`+userTable+`" where nickname=$1 or email=$2 limit 1`, nickname, email)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) GetUsersInDB(nickname string, email string) []model.User {
	var users []model.User
	err := r.db.Get(&users, `select * from "`+userTable+`" where nickname=$1 or email=$2`, nickname, email)
	checkErr(err)
	return users
}

func (r *DatabaseRepository) СreateUserInDB(nickname string, user model.User) model.User {
	user.Nickname = nickname
	_, err := r.db.Exec(`insert into "`+userTable+`" (nickname, email, about, fullname) values ($1, $2, $3, $4)`,
		user.Nickname, user.Email, user.About, user.Fullname)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) GetUserIdByName(nickname string) int32 {
	var userId int32
	err := r.db.Get(&userId, `select id from "`+userTable+`" where nickname=$1`, nickname)
	checkErr(err)
	return userId
}

func (r *DatabaseRepository) GetUserById(id int32) model.User {
	var user model.User
	err := r.db.Get(&user, `select * from "`+userTable+`" where id=$1`, id)
	checkErr(err)
	return user
}
package repository

import "TP_DB/model"

func (r *DatabaseRepository) IsUserInDB(nickname string, email string) bool {
	var count int
	if nickname != "" {
		err := r.db.Get(&count, `select count(*) from "`+userTable+`" where nickname=$1`, nickname)
		checkErr(err)
		if count != 0 {
			return true
		}
	}
	if email != "" {
		err := r.db.Get(&count, `select count(*) from "`+userTable+`" where email=$1`, email)
		checkErr(err)
		if count != 0 {
			return true
		}
	}
	return false
}

func (r *DatabaseRepository) GetUserInDB(nickname string, email string) model.User {
	var user, emptyUser model.User
	err := r.db.Get(&user, `select * from "`+userTable+`" where nickname=$1`, nickname)
	checkErr(err)
	if user == emptyUser {
		err := r.db.Get(&user, `select * from "`+userTable+`" where email=$1`, email)
		checkErr(err)
	}
	return user
}

func (r *DatabaseRepository) СreateUserInDB(nickname string, user model.User) model.User {
	user.Nickname = nickname
	_, err := r.db.Exec(`insert into "`+userTable+`" (nickname, email, about, fullname) values ($1, $2, $3, $4)`,
		user.Nickname, user.Email, user.About, user.Fullname)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) GetUserIdByName(nickname string) int32 {
	var userId int32 = -1
	err := r.db.Get(&userId, `select id from "`+userTable+`" where nickname=$1`, nickname)
	checkErr(err)
	return userId
}
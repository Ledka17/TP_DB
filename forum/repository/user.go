package repository

import (
	"github.com/Ledka17/TP_DB/model"
)

func (r *DatabaseRepository) IsUserInDB(nickname string, email string) bool {
	count := 0
	countMail := 0
	// Один запрос
	err := r.db.Get(&count, `select count(*) from "`+userTable+`" where lower(nickname)=lower($1) or lower(email)=lower($2)`, nickname, email)
	checkErr(err)
	count += countMail

	// Вернуть юзера
	if count != 0 {
		return true
	}
	return false
}

func (r *DatabaseRepository) IsUsersInDB(nicknames []string) bool {
	for _, nickname := range nicknames {
		if !r.IsUserInDB(nickname, "") {
			return false
		}
	}
	return true
}

func (r *DatabaseRepository) GetUserInDB(nickname string, args ...string) model.User {
	var user model.User
	email := ""
	if args != nil {
		email = args[0]
	}
	err := r.db.Get(&user, `select * from "`+userTable+`" where lower(nickname)=lower($1) or lower(email)=lower($2) limit 1`, nickname, email)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) GetUsersInDB(nickname string, email string) []*model.User {
	var users []*model.User
	err := r.db.Select(&users, `select * from "`+userTable+`" where lower(nickname)=lower($1) or lower(email)=lower($2)`, nickname, email)
	checkErr(err)
	return users
}

func (r *DatabaseRepository) CreateUserInDB(nickname string, user model.User) model.User {
	user.Nickname = nickname
	_, err := r.db.Exec(`insert into "`+userTable+`" (nickname, email, about, fullname) values ($1, $2, $3, $4)`,
		user.Nickname, user.Email, user.About, user.Fullname)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) GetUserIdByName(nickname string) int32 {
	var userId int32
	err := r.db.Get(&userId, `select id from "`+userTable+`" where lower(nickname)=lower($1)`, nickname)
	checkErr(err)
	return userId
}

func (r *DatabaseRepository) GetUserById(id int32) model.User {
	var user model.User
	err := r.db.Get(&user, `select * from "`+userTable+`" where id=$1`, id)
	checkErr(err)
	return user
}

func (r *DatabaseRepository) ChangeUserInDB(nickname string, userUpdate model.UserUpdate) model.User {
	if userUpdate.Fullname != "" {
		_, err := r.db.Exec(
			`update "`+userTable+`" set fullname=$1 where lower(nickname)=lower($2)`,
			userUpdate.Fullname, nickname,
		)
		checkErr(err)
	}
	if userUpdate.About != "" {
		_, err := r.db.Exec(
			`update "`+userTable+`" set about=$1 where lower(nickname)=lower($2)`,
			userUpdate.About, nickname,
		)
		checkErr(err)
	}
	if userUpdate.Email != "" {
		_, err := r.db.Exec(
			`update "`+userTable+`" set email=$1 where lower(nickname)=lower($2)`,
			userUpdate.Email, nickname,
		)
		checkErr(err)
	}
	user := r.GetUserInDB(nickname)
	return user
}
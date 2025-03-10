package users

import (
	"database/sql"
	"maxwarden/database"
)

type User struct {
	ID             int32  `db:"id"`
	Username       string `db:"username"`
	Email          string `db:"email"`
	Firstname      string `db:"firstname"`
	Lastname       string `db:"lastname"`
	Password       string `db:"password"`
	FailedAttempts int32  `db:"failed_attempts"`
	SecurityStamp  string `db:"security_stamp"`
	LastLogin      string `db:"last_login"`
	Data           []byte `db:"data"`
}

func FetchById(id int32) (User, error) {
	qb := &database.QueryBuilder{}
	qb.BaseSQL = "SELECT * FROM users u WHERE u.id = ?"

	return database.Get[User](qb, database.DB, id)
}

func FetchByUsername(username string) (User, error) {
	qb := &database.QueryBuilder{}
	qb.BaseSQL = "SELECT * FROM users u WHERE u.username = ?"

	return database.Get[User](qb, database.DB, username)
}

func FetchSecurityStamp(userid int) (string, error) {
	qb := &database.QueryBuilder{}
	qb.BaseSQL = "SELECT u.security_stamp FROM users u WHERE u.id = ?"

	return database.Get[string](qb, database.DB, userid)
}

func Update(user User) (sql.Result, error) {
	qb := &database.QueryBuilder{}
	qb.BaseSQL = "UPDATE users"

	qb.Setters = []database.QuerySetter{
		{Column: "failed_attempts", Parameter: user.FailedAttempts},
		{Column: "data", Parameter: user.Data},
	}

	qb.Where = []database.QueryFilter{
		{Column: "id", Parameter: user.ID},
	}

	return database.Update[User](qb, database.DB)
}

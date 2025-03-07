package users

import (
	"database/sql"
	"maxwarden/.jet/model"
	"maxwarden/database"

	. "maxwarden/.jet/table"

	. "github.com/go-jet/jet/v2/sqlite"
)

func FetchById(id int32) (model.User, error) {
	user := model.User{}

	stmt := SELECT(
		User.AllColumns,
	).FROM(User).WHERE(
		User.ID.EQ(Int32(int32(id))),
	)

	err := stmt.Query(database.DB, &user)

	return user, err
}

func FetchByUsername(username string) (model.User, error) {
	dest := model.User{}

	stmt := SELECT(
		User.AllColumns,
	).FROM(
		User,
	).WHERE(
		User.Username.EQ(String(username)),
	)

	err := stmt.Query(database.DB, &dest)
	return dest, err
}

func FetchSecurityStamp(userid int) (string, error) {
	stamp := ""

	stmt := SELECT(
		User.SecurityStamp,
	).FROM(User).WHERE(
		User.ID.EQ(Int32(int32(userid))),
	)

	err := stmt.Query(database.DB, &stamp)
	return stamp, err
}

func Update(user model.User) (sql.Result, error) {
	stmt := User.UPDATE(User.MutableColumns).
		MODEL(user).
		WHERE(User.ID.EQ(Int32(user.ID)))

	return stmt.Exec(database.DB)
}
package auth

import (
	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"

	db "github.com/aodin/listofthings/db"
)

type UserManager struct {
	conn sql.Connection
}

func (m *UserManager) Create(name, email string) db.User {
	// TODO prevent duplicate emails
	user := db.NewUser(name, email)
	stmt := pg.Insert(db.Users).Values(user).Returning(db.Users)
	m.conn.MustQueryOne(stmt, &user)
	return user
}

func (m *UserManager) Get(id int64) (user db.User) {
	stmt := db.Users.Select().Where(db.Users.C["id"].Equals(id))
	m.conn.MustQueryOne(stmt, &user)
	return
}

func Users(conn sql.Connection) *UserManager {
	return &UserManager{
		conn: conn,
	}
}

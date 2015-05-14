package db

import (
	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"

	"github.com/aodin/listofthings/db/fields"
)

type User struct {
	ID    int64  `db:"id,omitempty" json:"id,omitempty"`
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
	fields.Timestamp
}

func (user User) Exists() bool {
	return user.ID != 0
}

func (user User) String() string {
	if user.Name == "" {
		return "Anonymous User"
	}
	return string(user.Name)
}

func NewUser(name, email string) User {
	return User{Name: name, Email: email}
}

var Users = sql.Table("users",
	sql.Column("id", pg.Serial{NotNull: true}),
	sql.Column("email", sql.String{NotNull: true, Length: 256}),
	sql.Column("name", sql.String{Length: 128, NotNull: true}),
	sql.Column("created_at", sql.Timestamp{NotNull: true, Default: pg.Now}),
	sql.Column("updated_at", sql.Timestamp{}),
	sql.Column("deleted_at", sql.Timestamp{}),
	sql.PrimaryKey("id"),
)

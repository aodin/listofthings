package db

import (
	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"

	"github.com/aodin/listofthings/db/fields"
)

type User struct {
	ID    fields.ID    `json:"id,omitempty"`
	Email fields.Email `json:"email"`
	Name  fields.Name  `json:"name"`
	fields.Timestamp
}

func (user User) Exists() bool {
	return user.ID.Exists()
}

func (user User) String() string {
	if user.Name == "" {
		return "Anonymous User"
	}
	return string(user.Name)
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

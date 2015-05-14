package db

import (
	"time"

	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"
)

type Session struct {
	Key     string    `db:"key"`
	UserID  int64     `db:"user_id"`
	Expires time.Time `db:"expires_at"`
}

func (session Session) Exists() bool {
	return session.Key != ""
}

var Sessions = sql.Table("sessions",
	sql.Column("key", sql.String{NotNull: true}),
	sql.ForeignKey(
		"user_id",
		Users.C["id"],
		sql.Integer{NotNull: true},
	).OnDelete(sql.Cascade),
	sql.Column("expires_at", sql.Timestamp{NotNull: true}),
	sql.PrimaryKey("key"),
)

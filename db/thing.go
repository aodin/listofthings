package db

import (
	"fmt"

	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"

	"github.com/aodin/listofthings/db/fields"
)

const MaxNameLength = 256

// Thing is a thing with a name
type Thing struct {
	ID      int64  `db:"id,omitempty" json:"id"`
	Name    string `db:"-" json:"name"`
	Content string `db:"content" json:"-"`
	fields.Timestamp
}

func (t Thing) String() string {
	return t.Name
}

func (t Thing) Error() error {
	if t.Name == "" {
		return fmt.Errorf("Names cannot be blank")
	}
	if len(t.Name) > MaxNameLength {
		return fmt.Errorf(
			"Names cannot be longer than %d characters",
			MaxNameLength,
		)
	}
	return nil
}

func (t Thing) Values() sql.Values {
	return sql.Values{
		"content": t.Content,
	}
}

func NewThing(name string) Thing {
	return Thing{Name: name}
}

var Things = sql.Table("things",
	sql.Column("id", pg.Serial{NotNull: true}),
	sql.Column("content", pg.JSON{NotNull: true}),
	sql.Column("created_at", sql.Timestamp{NotNull: true, Default: pg.Now}),
	sql.Column("updated_at", sql.Timestamp{}),
	sql.Column("deleted_at", sql.Timestamp{}),
	sql.PrimaryKey("id"),
)

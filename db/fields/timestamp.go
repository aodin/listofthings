package fields

import "time"

type Timestamp struct {
	CreatedAt time.Time  `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (t Timestamp) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t Timestamp) Age() time.Duration {
	return t.age(time.Now().UTC())
}

func (t Timestamp) age(now time.Time) time.Duration {
	return now.Sub(t.CreatedAt)
}

func (t Timestamp) LastActivity() time.Time {
	if t.DeletedAt != nil {
		return *t.DeletedAt
	}
	if t.UpdatedAt != nil {
		return *t.UpdatedAt
	}
	return t.CreatedAt
}

func NewTimestamp() Timestamp {
	return newTimestamp(time.Now().UTC())
}

func newTimestamp(now time.Time) Timestamp {
	return Timestamp{CreatedAt: now}
}

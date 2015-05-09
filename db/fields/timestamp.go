package fields

import "time"

type Timestamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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

func NewTimestamp() Timestamp {
	return newTimestamp(time.Now().UTC())
}

func newTimestamp(now time.Time) Timestamp {
	return Timestamp{CreatedAt: now}
}

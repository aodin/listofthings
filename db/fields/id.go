package fields

type ID int64

func (id ID) Exists() bool {
	return id != 0
}

package postgres

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aodin/aspect"
)

type Serial struct {
	PrimaryKey bool
	NotNull    bool
}

var _ aspect.Type = Serial{}

func (s Serial) Create(d aspect.Dialect) (string, error) {
	compiled := "SERIAL"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Serial) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Serial) IsRequired() bool {
	return s.NotNull
}

func (s Serial) IsUnique() bool {
	return true
}

func (s Serial) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float64:
		v := int64(t)
		if t != float64(v) {
			return i, fmt.Errorf(
				"value is numeric, but not a whole number: %f",
				t,
			)
		}
		i = v
	case int:
		i = int64(i.(int))
	case int64:
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}

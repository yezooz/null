// Package null contains types that consider zero input and null input as separate values.
// Types in this package will always encode to their null value if null.
// Use the zero subpackage if you want empty and null to be treated the same.
package null

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"reflect"
	"time"
)

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type Time struct {
	pq.NullTime
}

// TimeFrom creates a new Time that will never be blank.
func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

// TimeFromPtr creates a new Time that be null if s is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Now(), false)
	}
	return NewTime(*t, true)
}

// NewTime creates a new Time
func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: pq.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
// It also supports unmarshalling a sql.NullString.
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch x := v.(type) {
	case time.Time:
		t.Time = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &t.NullTime)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil && !t.Time.IsZero()
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Time if the input is a blank string.
func (t *Time) UnmarshalText(text []byte) error {
	t.Time = t.Time
	t.Valid = !t.Time.IsZero()
	return nil
}

// SetValid changes this Time's value and also sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for null or empty strings, for future omitempty support. (Go 1.4?)
// Will return false s if blank but non-null.
func (t Time) IsZero() bool {
	return !t.Valid
}

package dot

import (
	"strconv"
	"time"
)

type ID int64

func ParseID(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	return ID(id), err
}

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type NullBool struct {
	Bool  bool
	Valid bool
}

type NullString struct {
	String string
	Valid  bool
}

type NullInt64 struct {
	Int64 int64
	Valid bool
}

type NullInt32 struct {
	Int32 int32
	Valid bool
}

type NullInt struct {
	Int   int
	Valid bool
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

type NullID struct {
	ID    ID
	Valid bool
}

func (id ID) Wrap() NullID          { return NullID{ID: id, Valid: true} }
func WrapID(id ID) NullID           { return NullID{ID: id, Valid: true} }
func Bool(b bool) NullBool          { return NullBool{Bool: b, Valid: true} }
func String(s string) NullString    { return NullString{String: s, Valid: true} }
func Int64(i int64) NullInt64       { return NullInt64{Int64: i, Valid: true} }
func Int32(i int32) NullInt32       { return NullInt32{Int32: i, Valid: true} }
func Int(i int) NullInt             { return NullInt{Int: i, Valid: true} }
func Float64(f float64) NullFloat64 { return NullFloat64{Float64: f, Valid: true} }

func PID(i *ID) NullID {
	if i == nil {
		return NullID{}
	}
	return NullID{ID: *i, Valid: true}
}

func (n NullString) Apply(s string) string {
	if n.Valid {
		return n.String
	}
	return s
}

func (n NullBool) Apply(b bool) bool {
	if n.Valid {
		return n.Bool
	}
	return b
}

func (n NullInt64) Apply(i int64) int64 {
	if n.Valid {
		return n.Int64
	}
	return i
}

func (n NullInt32) Apply(i int32) int32 {
	if n.Valid {
		return n.Int32
	}
	return i
}

func (n NullInt) Apply(i int) int {
	if n.Valid {
		return n.Int
	}
	return i
}

func (id NullID) Apply(i ID) ID {
	if id.Valid {
		return id.ID
	}
	return i
}

var zeroTime = time.Unix(0, 0)

// IsZeroTime checks whether given time is zero.
// When transport time using grpc, empty time is marshalled to time.Unix(0, 0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(zeroTime)
}

type Time time.Time

func (t Time) String() string {
	return t.ToTime().Format(time.RFC3339)
}

func (t Time) ToTime() time.Time {
	return time.Time(t)
}

func (t Time) IsZero() bool {
	return time.Time(t).IsZero() || time.Time(t).Equal(zeroTime)
}

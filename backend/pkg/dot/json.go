package dot

import (
	"strconv"
	"time"
)

const null = "null"

var jsonNull = []byte("null")
var jsonTrue = []byte("true")
var jsonFalse = []byte("false")

func (id ID) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 32)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(id), 10)
	b = append(b, '"')
	return b, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*id = 0
		return nil
	}
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		*id = 0
		return err
	}
	*id = ID(i)
	return err
}

func (id NullID) MarshalJSON() ([]byte, error) {
	if !id.Valid {
		return jsonNull, nil
	}
	return id.ID.MarshalJSON()
}

func (id *NullID) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*id = NullID{}
		return nil
	}
	var _id ID
	err := _id.UnmarshalJSON(data)
	if err != nil {
		*id = NullID{}
		return err
	}
	*id = WrapID(_id)
	return nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	if n.Bool {
		return jsonTrue, nil
	}
	return jsonFalse, nil
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullBool{}
		return nil
	}
	b, err := strconv.ParseBool(string(data))
	if err != nil {
		*n = NullBool{}
		return err
	}
	*n = Bool(b)
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	dst := make([]byte, 0, len(n.String)+2+len(n.String)/4)
	return strconv.AppendQuote(dst, n.String), nil
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullString{}
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		*n = NullString{}
	}
	*n = String(s)
	return nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	b := strconv.AppendInt(nil, n.Int64, 10)
	return b, nil
}

func (n *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullInt64{}
		return nil
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		*n = NullInt64{}
	}
	*n = Int64(i)
	return nil
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	b := strconv.AppendInt(nil, int64(n.Int32), 10)
	return b, nil
}

func (n *NullInt32) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullInt32{}
		return nil
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		*n = NullInt32{}
	}
	*n = Int32(int32(i))
	return nil
}

func (n NullInt) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	b := strconv.AppendInt(nil, int64(n.Int), 10)
	return b, nil
}

func (n *NullInt) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullInt{}
		return nil
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		*n = NullInt{}
	}
	*n = Int(int(i))
	return nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return jsonNull, nil
	}
	b := strconv.AppendFloat(nil, n.Float64, 'f', -1, 64)
	return b, nil
}

func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == null {
		*n = NullFloat64{}
		return nil
	}
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		*n = NullFloat64{}
	}
	*n = Float64(f)
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if IsZeroTime(time.Time(t)) {
		return jsonNull, nil
	}
	tt := time.Time(t)
	tt = tt.Add(-time.Duration(tt.Nanosecond()))
	return tt.MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == null || string(data) == `""` {
		*t = Time{}
		return nil
	}
	var tt time.Time
	err := tt.UnmarshalJSON(data)
	if err != nil || IsZeroTime(tt) {
		*t = Time{}
	} else {
		*t = Time(tt)
	}
	return err
}

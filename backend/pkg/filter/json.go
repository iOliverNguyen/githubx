package filter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ng-vu/githubx/backend/pkg/dot"
	"github.com/ng-vu/githubx/backend/pkg/xerrors"
)

const null = "null"

var bytesComma = []byte{','}

func (ids *IDs) UnmarshalJSON(data []byte) error {
	if string(data) == null || string(data) == `""` {
		*ids = nil
		return nil
	}
	if data[0] == '[' {
		var _ids []dot.ID
		if err := json.Unmarshal(data, &_ids); err != nil {
			return xerrors.Errorf(err, "json: can not read `%s` as id", data)
		}
		if len(_ids) == 0 {
			return xerrors.Errorf(nil, "json: can not read `%s` as id", data)
		}
		*ids = _ids
		return nil
	}
	parts := bytes.Split(unsafeUnquote(data), bytesComma)
	_ids := make(IDs, len(parts))
	for i, p := range parts {
		if err := json.Unmarshal(p, &_ids[i]); err != nil {
			return xerrors.Errorf(err, "json: can not read `%s` as id", data)
		}
		if _ids[i] == 0 {
			return xerrors.Errorf(nil, "json: can not read `%s` as id", data)
		}
	}
	if len(_ids) == 0 {
		return xerrors.Errorf(nil, "json: can not read `%s` as id", data)
	}
	*ids = _ids
	return nil
}

func (ss *Strings) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*ss = nil
		return nil
	}
	if data[0] == '[' {
		if err := json.Unmarshal(data, (*[]string)(ss)); err != nil {
			return xerrors.Errorf(err, "json: can not read `%s` as string array", data)
		}
		if len(*ss) == 0 {
			return xerrors.Errorf(nil, "json: can not read `%s` as string array", data)
		}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ss = []string{s}
	return nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*d = Date{}
		return nil
	}
	if data[0] == '{' {
		if err := json.Unmarshal(data, (*date)(d)); err != nil {
			return xerrors.Errorf(err, "json: can not read `%s` as date", data)
		}
		return nil
	}

	var dates []dot.Time // for correctly handle zero time
	parts := bytes.Split(unsafeUnquote(data), bytesComma)
	dates = make([]dot.Time, len(parts))
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		t, err := time.Parse(time.RFC3339, string(p))
		if err != nil {
			return xerrors.Errorf(err, "json: can not read `%s` as date", data)
		}
		dates[i] = dot.Time(t)
	}
	switch len(dates) {
	case 1:
		d.From = dates[0]
		d.To = dot.Time{}

	case 2:
		d.From = dates[0]
		d.To = dates[1]
		if !d.To.IsZero() &&
			d.From.ToTime().After(d.To.ToTime()) {
			return xerrors.Errorf(nil, "date range [%v,%v) is invalid", d.From, d.To)
		}

	default:
		return xerrors.Errorf(nil, "json: can not read `%s` as date", data)
	}

	if d.From.IsZero() && d.To.IsZero() {
		return xerrors.Errorf(nil, "json: can not read `%s` as date", data)
	}
	return nil
}

func unsafeUnquote(data []byte) []byte {
	if data[0] == '"' {
		return data[1 : len(data)-1]
	}
	return data
}

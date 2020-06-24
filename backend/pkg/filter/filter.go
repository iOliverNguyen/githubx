package filter

import (
	"github.com/ng-vu/githubx/backend/pkg/dot"
)

// Filter by id value
type IDs []dot.ID

// Filter by string value
type Strings []string

// Filter by date value
//
// Object:
//
//     - {"from":"2020-01-01T00:00:00Z","to":"2020-02-01T00:00:00Z"}
//
// String:
//
//     - "2020-01-01T00:00:00Z,2010-02-01T00:00:00Z":
//     - "2020-01-01T00:00:00Z,":
//     - ",2020-01-01T00:00:00Z":
type Date date

type date struct {
	From dot.Time `json:"from"`
	To   dot.Time `json:"to"`
}

func (d Date) IsZero() bool {
	return d.From.IsZero() && d.To.IsZero()
}

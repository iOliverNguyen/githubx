package store

import (
	"github.com/tidwall/buntdb"

	"github.com/ng-vu/githubx/backend/pkg/graph"
)

func (s *Datastore) SaveComments(issueNumber int, cmts ...*graph.Comment) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		for _, cmt := range cmts {
			if err := saveComment(tx, issueNumber, cmt); err != nil {
				return err
			}
		}
		return nil
	})
}

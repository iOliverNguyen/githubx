package store

import (
	"github.com/tidwall/buntdb"

	"github.com/ng-vu/githubx/backend/pkg/api"
	"github.com/ng-vu/githubx/backend/pkg/graph"
)

func (s *Datastore) SaveIssues(issues ...*graph.Issue) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		for _, is := range issues {
			if err := saveIssue(tx, is); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Datastore) LoadIssue(number interface{}) (is *api.Issue, err error) {
	err = s.db.View(func(tx *buntdb.Tx) error {
		is, err = loadIssue(tx, number)
		return err
	})
	return is, err
}

func (s *Datastore) ListIssues(limit int) (_ []*api.Issue, err error) {
	result := make([]*api.Issue, 0, limit)
	err = s.db.View(func(tx *buntdb.Tx) error {
		pattern := keyIssue(nil, "")
		return tx.AscendKeys(pattern, func(key, value string) bool {
			is, err2 := loadIssue(tx, extractID(key, 1))
			if err2 != nil {
				err = err2
				return false
			}
			result = append(result, is)
			return len(result) < limit
		})
	})
	return result, err
}

func (s *Datastore) ListIssuesByIndex(index string, limit int) (_ []*api.Issue, err error) {
	result := make([]*api.Issue, 0, limit)
	err = s.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(index, func(key, value string) bool {
			is, err2 := loadIssue(tx, extractID(key, 1))
			if err2 != nil {
				err = err2
				return false
			}
			result = append(result, is)
			return len(result) < limit
		})
	})
	return result, err
}

func (s *Datastore) ListIssuesByIndexDesc(index string, limit int) (_ []*api.Issue, err error) {
	result := make([]*api.Issue, 0, limit)
	err = s.db.View(func(tx *buntdb.Tx) error {
		return tx.Descend(index, func(key, value string) bool {
			is, err2 := loadIssue(tx, extractID(key, 1))
			if err2 != nil {
				err = err2
				return false
			}
			result = append(result, is)
			return len(result) < limit
		})
	})
	return result, err
}

/* Package store stores and retrieves data.

Key structure:

    is:number:model        issue
    is:number:lastcmt      last comment
    is:number:labels       label ids
	is:number:cmt:id:model comment
    label:id:model         label
*/
package store

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/buntdb"

	"github.com/ng-vu/githubx/backend/pkg/api"
	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/graph"
	"github.com/ng-vu/githubx/backend/pkg/model"
)

const (
	keyLastCmt       = "lastCmt"
	keyLastChangedAt = "lastChangedAt"
	keyLabels        = "labels"
	keyModel         = "model"

	IndexIssueCreatedAt     = "idxIssueCreatedAt"
	IndexIssueUpdatedAt     = "idxIssueUpdatedAt"
	IndexIssueLastChangedAt = "idxIssueLastChangedAt"
	IndexCmtCreatedAt       = "idxCmtCreatedAt"
	IndexCmtIssueCreatedAt  = "idxCmtIssueCreatedAt"
)

type Issue struct {
	model.Issue

	LastComment string
}

// 000:001:002:003:004
func extractID(key string, pos int) string {
	parts := strings.SplitN(key, ":", pos+2)
	return parts[pos]
}

func ids(ids ...string) string   { return strings.Join(ids, ",") }
func parseIDs(s string) []string { return strings.Split(s, ",") }

func appendNumber(b []byte, number interface{}) []byte {
	switch number := number.(type) {
	case string:
		b = append(b, number...)
	case int:
		b = strconv.AppendInt(b, int64(number), 10)
	default:
		if number == nil {
			b = append(b, '*')
		} else {
			log.Printf("%#v (%T)", number, number)
			panic("unexpected number")
		}
	}
	return b
}

func keyIssue(number interface{}, subkey string) string {
	if subkey == "" {
		subkey = "model"
	}
	length := 8 + len(subkey)
	b := make([]byte, 0, length)
	b = append(b, "is:"...)
	b = appendNumber(b, number)
	b = append(b, ':')
	b = append(b, subkey...)
	return unsafeBytesToString(b)
}

func keyComment(number interface{}, id string, subkey string) string {
	if subkey == "" {
		subkey = "model"
	}
	b := make([]byte, 0, 8+len(id)+len(subkey)+4)
	b = append(b, "is:"...)
	b = appendNumber(b, number)
	b = append(b, "cmt:"...)
	if id == "" {
		b = append(b, '*')
	} else {
		b = append(b, id...)
	}
	b = append(b, ':')
	b = append(b, subkey...)
	return unsafeBytesToString(b)
}

func saveIssue(tx *buntdb.Tx, is *graph.Issue) error {
	if is.Body != "" && (is.BodyHTML == "" || is.BodyText == "") {
		var err error
		is.BodyHTML, is.BodyText, err = github.ConvertMD(is.Body)
		if err != nil {
			log.Println("ConvertMD", err)
		}
	}
	if is.Number <= 0 {
		log.Printf("%#v", is)
		panic("invalid number")
	}
	_, _, err := tx.Set(keyIssue(is.Number, ""), jsonMarshal(is), nil)
	if err != nil {
		return err
	}

	lastChangedAt := is.UpdatedAt
	if lastCmt := is.Comments.Last(); lastCmt != nil && lastCmt.UpdatedAt.After(lastChangedAt) {
		lastChangedAt = lastCmt.UpdatedAt
	}
	if err = saveIssueLastChangedAt(tx, is.Number, lastChangedAt); err != nil {
		return err
	}

	if err = saveComments(tx, is.Number, is.Comments); err != nil {
		return err
	}
	return nil
}

func loadIssue(tx *buntdb.Tx, number interface{}) (*api.Issue, error) {
	var is api.Issue
	val, _ := tx.Get(keyIssue(number, ""))
	if val == "" {
		return nil, nil
	}
	if err := jsonUnmarshal(val, &is.Issue); err != nil {
		return nil, err
	}

	is.LastChangedAt = loadIssueLastChangedAt(tx, number)
	cmts, err := loadCommentsForIssue(tx, number)
	if err != nil {
		return nil, err
	}
	is.Comments = cmts
	return &is, nil
}

func saveComments(tx *buntdb.Tx, issueNumber int, cmts *graph.Comments) error {
	if cmts == nil {
		return nil
	}
	for _, cmt := range cmts.Nodes {
		if err := saveComment(tx, issueNumber, cmt); err != nil {
			return err
		}
	}
	return nil
}

func saveComment(tx *buntdb.Tx, issueNumber int, cmt *graph.Comment) error {
	if cmt.Body != "" && (cmt.BodyHTML == "" || cmt.BodyText == "") {
		var err error
		cmt.BodyHTML, cmt.BodyText, err = github.ConvertMD(cmt.Body)
		if err != nil {
			log.Println("ConvertMD", err)
		}
	}
	if issueNumber <= 0 || cmt.ID == "" {
		log.Printf("%v %v", issueNumber, cmt.ID)
		panic("invalid number or id")
	}
	_, _, err := tx.Set(keyComment(issueNumber, cmt.ID, ""), jsonMarshal(cmt.Comment), nil)
	if err != nil {
		return err
	}
	return saveIssueLastChangedAt(tx, issueNumber, cmt.UpdatedAt)
}

func loadComment(tx *buntdb.Tx, issueNumber int, id string) (*api.Comment, error) {
	var cmt api.Comment
	val, _ := tx.Get(keyComment(issueNumber, id, ""))
	if val == "" {
		return nil, nil
	}
	if err := jsonUnmarshal(val, &cmt.Comment); err != nil {
		return nil, err
	}
	return &cmt, nil
}

func loadCommentsForIssue(tx *buntdb.Tx, number interface{}) ([]*api.Comment, error) {
	result := make([]*api.Comment, 0, 32)
	pattern := keyComment(number, "", "")
	err := tx.AscendKeys(pattern, func(key, value string) bool {
		var cmt api.Comment
		if err := jsonUnmarshal(value, &cmt); err != nil {
			return false
		}
		result = append(result, &cmt)
		return true
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.Before(result[j].UpdatedAt)
	})
	return result, err
}

func saveIssueLastChangedAt(tx *buntdb.Tx, number interface{}, t time.Time) (err error) {
	kLastChangedAt := keyIssue(number, keyLastChangedAt)
	currentTime := parseTime(txGet(tx, kLastChangedAt))
	if t.After(currentTime) {
		_, _, err = tx.Set(kLastChangedAt, formatTime(t), nil)
	}
	return
}

func loadIssueLastChangedAt(tx *buntdb.Tx, number interface{}) time.Time {
	kLastChangedAt := keyIssue(number, keyLastChangedAt)
	return parseTime(txGet(tx, kLastChangedAt))
}

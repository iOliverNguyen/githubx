package api

import (
	"time"

	"github.com/ng-vu/githubx/backend/pkg/filter"
	"github.com/ng-vu/githubx/backend/pkg/model"
)

type AuthorizeRequest struct {
}

type AuthorizeResponse struct {
	Username string `json:"username"`
}

type Issue struct {
	model.Issue

	Assignees []*model.User

	Labels []*model.Label

	Milestone *model.Milestone

	Comments []*Comment `json:"comments"`

	// this field is the latest of LastComment.UpdatedAt and Issue.UpdatedAt
	LastChangedAt time.Time `json:"lastChangedAt"`
}

type Comment struct {
	model.Comment
}

type ListIssuesRequest struct {
	Number  filter.IDs     `json:"id"`
	Project filter.Strings `json:"project"`
	Label   filter.Strings `json:"label"`

	OrderBy string `json:"orderBy"`
	Limit   int    `json:"limit"`
}

type ListIssuesResponse struct {
	Issues []*Issue `json:"issues"`
}

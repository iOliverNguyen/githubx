package graph

import "github.com/ng-vu/githubx/backend/pkg/model"

type ListIssuesResponse struct {
	Organization *Organization `json:"organization"`
}

type Organization struct {
	Respository *Repository `json:"repository"`
}

type Repository struct {
	Issues *Issues `json:"issues"`
}

type Issues struct {
	Nodes      []*Issue  `json:"nodes"`
	TotalCount int       `json:"totalCount"`
	PageInfo   *PageInfo `json:"pageInfo"`
}

type Issue struct {
	model.Issue

	Assignees *Assignees `json:"assignees"`
	Author    *Author    `json:"author"`
	Comments  *Comments  `json:"comments"`
	Labels    *Labels    `json:"labels"`
}

type Assignees struct {
	Nodes      []*Assignee `json:"nodes"`
	TotalCount int         `json:"totalCount"`
}

type Assignee = model.User

type Author struct {
	model.Author
}

type Comments struct {
	Nodes      []*Comment `json:"nodes"`
	TotalCount int        `json:"totalCount"`
	PageInfo   *PageInfo  `json:"pageInfo"`
}

func (cs *Comments) Last() *Comment {
	if cs == nil {
		return nil
	}
	if len(cs.Nodes) == 0 {
		return nil
	}
	if cs.Nodes[0].CreatedAt.After(cs.Nodes[len(cs.Nodes)-1].CreatedAt) {
		return cs.Nodes[0]
	}
	return cs.Nodes[len(cs.Nodes)-1]
}

type Comment struct {
	model.Comment

	Author *Author `json:"author"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	HasPrevPage bool   `json:"hasPreviousPage"`
	EndCursor   string `json:"endCursor"`
	StartCursor string `json:"startCursor"`
}

type Labels struct {
	Nodes      []*Label `json:"nodes"`
	TotalCount int      `json:"totalCount"`
}

type Label = model.Label

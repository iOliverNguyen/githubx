package github

import (
	"fmt"
	"strings"

	"github.com/ng-vu/githubx/backend/pkg/graph"
)

func pagingAfter(limit int, after string) string {
	if after == "" {
		return fmt.Sprintf("first: %v", limit)
	}
	return fmt.Sprintf("first: %v, after: %q", limit, after)
}

func LoginQuery(arg OrgRepo) string {
	var b strings.Builder
	must(loginTpl.Execute(&b, arg))
	return b.String()
}

type LoginResponse struct {
	Viewer struct {
		ID    string `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	} `json:"viewer"`
}

type OrgRepo struct {
	Org  string `yaml:"org"`
	Repo string `yaml:"repo"`
}

func (r OrgRepo) FullName() string { return r.Org + "/" + r.Repo }

type ListIssueQueryArgs struct {
	OrgRepo
}

func ListIssuesQuery(arg ListIssueQueryArgs, after string) string {
	var b strings.Builder
	must(listIssueTpl.Execute(&b, M{
		"Org":         arg.Org,
		"Repo":        arg.Repo,
		"IssuePaging": pagingAfter(100, after),
	}))
	return b.String()
}

type ListIssuesResponse struct {
	Organization *graph.Organization `json:"organization"`
}

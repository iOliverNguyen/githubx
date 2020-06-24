package loader

import (
	"context"
	"log"

	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/store"
)

type Loader struct {
	store   *store.Datastore
	client  *github.Client
	orgRepo github.OrgRepo
}

func NewLoader(client *github.Client, st *store.Datastore, cfg github.OrgRepo) *Loader {
	return &Loader{
		store:   st,
		client:  client,
		orgRepo: cfg,
	}
}

func (ld *Loader) LoadAllIssues(ctx context.Context) error {
	arg := github.ListIssueQueryArgs{OrgRepo: ld.orgRepo}
	count := 0
	next, nextCursor := true, ""
	for next {
		var resp github.ListIssuesResponse
		query := github.ListIssuesQuery(arg, nextCursor)
		_, err := ld.client.SendRequest(ctx, query, &resp)
		if err != nil {
			return err
		}
		issues := resp.Organization.Respository.Issues
		if err = ld.store.SaveIssues(issues.Nodes...); err != nil {
			return err
		}

		next, nextCursor = issues.PageInfo.HasNextPage, issues.PageInfo.EndCursor
		count += len(issues.Nodes)
		log.Printf("loaded %v/%v issues\n", count, issues.TotalCount)
	}
	return nil
}

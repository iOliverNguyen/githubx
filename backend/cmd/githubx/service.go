package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/ng-vu/githubx/backend/pkg/api"
	"github.com/ng-vu/githubx/backend/pkg/eventstream"
	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/session"
	"github.com/ng-vu/githubx/backend/pkg/store"
	"github.com/ng-vu/githubx/backend/pkg/sutil"
)

type Service struct {
	mu sync.RWMutex

	OrgRepo   github.OrgRepo
	Sessions  *session.Manager
	Client    *github.Client
	Store     *store.Datastore
	Publisher eventstream.Publisher
}

func NewService(
	orgRepo github.OrgRepo,
	client *github.Client,
	st *store.Datastore,
	sm *session.Manager,
	p eventstream.Publisher,
) *Service {
	s := &Service{
		OrgRepo:   orgRepo,
		Sessions:  sm,
		Client:    client,
		Store:     st,
		Publisher: p,
	}
	return s
}

func (s *Service) ServeAuthorize(w http.ResponseWriter, r *http.Request) {
	var req api.AuthorizeRequest
	if err := sutil.Request(r).Decode(&req); err != nil {
		sutil.Respond(w).Error(400, err)
		return
	}
	sutil.Respond(w).Resp(s.Authorize(r.Context(), r, &req))
}

func (s *Service) ServeListIssues(w http.ResponseWriter, r *http.Request) {
	var req api.ListIssuesRequest
	if err := sutil.Request(r).Decode(&req); err != nil {
		sutil.Respond(w).Error(400, err)
		return
	}
	sutil.Respond(w).Resp(s.ListIssues(r.Context(), r, &req))
}

func (s *Service) Authorize(ctx context.Context, r *http.Request, req *api.AuthorizeRequest) (*api.AuthorizeResponse, error) {
	token, err := session.GetBearer(r)
	if err != nil {
		return nil, err
	}
	ss, err := s.Sessions.Authorize(ctx, token)
	if err != nil {
		return nil, err
	}
	resp := &api.AuthorizeResponse{Username: ss.Login}
	return resp, err
}

func (s *Service) ListIssues(ctx context.Context, r *http.Request, req *api.ListIssuesRequest) (*api.ListIssuesResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 100
	}
	issues, err := s.Store.ListIssuesByIndexDesc(store.IndexIssueLastChangedAt, req.Limit)
	resp := &api.ListIssuesResponse{
		Issues: issues,
	}
	return resp, err
}

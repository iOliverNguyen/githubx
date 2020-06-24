package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/ng-vu/githubx/backend/pkg/eventstream"

	"github.com/ng-vu/githubx/backend/pkg/graph"
	"github.com/ng-vu/githubx/backend/pkg/model"
	"github.com/ng-vu/githubx/backend/pkg/sutil"
	"github.com/ng-vu/githubx/backend/pkg/webhook"
)

type M map[string]interface{}

const secret = "1@3$5^7*"

func (s *Service) ServeWebhook(w http.ResponseWriter, r *http.Request) {
	var req webhook.WebhookRequest
	if err := sutil.Request(r).Decode(&req); err != nil {
		sutil.Respond(w).Error(400, err)
		return
	}
	if req.Repository == nil || req.Repository.FullName != s.OrgRepo.FullName() {
		sutil.Respond(w).Resp(M{"code": "ignore unknown repo"}, nil)
		return
	}
	switch r.Header.Get("X-GitHub-Event") {
	case "":
		sutil.Respond(w).Error(400, errors.New("bad request"))

	case "ping":
		sutil.Respond(w).Resp(M{"code": "ok"}, nil)

	case "issues":
		sutil.Respond(w).Resp(s.webhookIssues(r.Context(), &req))

	case "issue_comment":
		sutil.Respond(w).Resp(s.webhookIssueComment(r.Context(), &req))

	case "project_card":
		sutil.Respond(w).Resp(s.webhookProjectCard(r.Context(), &req))

	default:
		sutil.Respond(w).Resp(M{"code": "ignored"}, nil)
	}
}

func convertIssue(is *webhook.Issue) *graph.Issue {
	_is := &graph.Issue{
		Issue:     model.Issue(*is.IssueX),
		Assignees: &graph.Assignees{Nodes: is.Assignees.ToUsers()},
		Author:    &graph.Author{Author: is.User.ToUser().ToAuthor()},
		Comments:  nil,
		Labels:    &graph.Labels{Nodes: is.Labels.ToLabels()},
	}
	return _is
}

func (s *Service) webhookIssues(ctx context.Context, req *webhook.WebhookRequest) (M, error) {
	is := convertIssue(req.Issue)
	if err := s.Store.SaveIssues(is); err != nil {
		return nil, err
	}
	_is, err := s.Store.LoadIssue(is.Number)
	if err != nil {
		return nil, err
	}
	e := eventstream.Event{
		Type:    "issue",
		Global:  true,
		Payload: M{"type": "issue", "data": _is},
	}
	s.Publisher.Publish(e)
	return M{"code": "ok"}, nil
}

func (s *Service) webhookIssueComment(ctx context.Context, req *webhook.WebhookRequest) (M, error) {
	is := convertIssue(req.Issue)
	if err := s.Store.SaveIssues(is); err != nil {
		return nil, err
	}
	cmt := &graph.Comment{
		Comment: (model.Comment)(*req.Comment.CommentX),
		Author:  &graph.Author{Author: req.Comment.User.ToUser().ToAuthor()},
	}
	if err := s.Store.SaveComments(is.Number, cmt); err != nil {
		return nil, err
	}
	_is, err := s.Store.LoadIssue(is.Number)
	if err != nil {
		return nil, err
	}
	e := eventstream.Event{
		Type:    "issue",
		Global:  true,
		Payload: M{"type": "issue", "data": _is},
	}
	s.Publisher.Publish(e)
	return M{"code": "ok"}, nil
}

func (s *Service) webhookProjectCard(ctx context.Context, req *webhook.WebhookRequest) (M, error) {
	return M{"code": "ok"}, nil
}

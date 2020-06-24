package session

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/xerrors"
)

type Session struct {
	Login string
	ID    string
}

type Manager struct {
	tokens  map[string]*Session
	client  *github.Client
	orgRepo github.OrgRepo
	mu      sync.RWMutex
}

func New(orgRepo github.OrgRepo, client *github.Client) *Manager {
	m := &Manager{
		tokens:  map[string]*Session{},
		client:  client,
		orgRepo: orgRepo,
	}
	return m
}

func (s *Manager) Authorize(ctx context.Context, token string) (*Session, error) {
	s.mu.RLock()
	if ss := s.Get(token); ss != nil {
		s.mu.RUnlock()
		return ss, nil
	}
	s.mu.RUnlock()

	var resp github.LoginResponse
	_, err := s.client.WithToken(token).SendRequest(ctx, github.LoginQuery(s.orgRepo), &resp)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	ss := Session{
		Login: resp.Viewer.Login,
		ID:    resp.Viewer.ID,
	}
	s.Set(token, ss)
	return &ss, nil
}

func (sm *Manager) Set(token string, s Session) {
	if s.Login == "" || s.ID == "" {
		panic("invalid session")
	}
	sm.tokens[token] = &s
}

func (sm *Manager) Get(token string) *Session {
	s := sm.tokens[token]
	if s == nil {
		return nil
	}
	ss := *s
	return &ss
}

func GetBearer(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	if token == "" {
		return "", xerrors.Errorf(nil, "unauthenticated")
	}
	return token, nil
}

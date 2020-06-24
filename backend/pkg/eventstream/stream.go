package eventstream

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/session"
	"github.com/ng-vu/githubx/backend/pkg/sutil"
)

type Publisher interface {
	Publish(event Event)
}

type Event struct {
	Type      string
	Global    bool // send to all users
	UserLogin string
	Payload   interface{}

	retryInSecond int
}

type Subscriber struct {
	ID int64

	UserLogin string

	AllEvents bool
	Events    []string

	ch chan *Event
}

type ToDelete struct {
	ID int64
	T  time.Time
}

type EventStream struct {
	Sessions     *session.Manager
	client       *github.Client
	subscribers  map[int64]*Subscriber
	eventChannel chan *Event
	toDeletes    []ToDelete

	ctx context.Context
	m   sync.RWMutex
}

func New(ctx context.Context, sm *session.Manager) *EventStream {
	es := &EventStream{
		Sessions:     sm,
		subscribers:  make(map[int64]*Subscriber),
		eventChannel: make(chan *Event, 256),
		ctx:          ctx,
	}
	go es.RunForwarder()
	go es.RunCleaner()
	return es
}

func (s *EventStream) Publish(event Event) {
	s.eventChannel <- &event
}

func (s *EventStream) RunForwarder() {
	for event := range s.eventChannel {
		s.forward(event)
	}
}

func (s *EventStream) RunCleaner() {
	for range time.Tick(10 * time.Second) {
		t := time.Now()
		s.m.RLock()
		if len(s.toDeletes) == 0 {
			s.m.RUnlock()
			continue
		}
		idx := 0
		for i, d := range s.toDeletes {
			if d.T.After(t) {
				idx = i
				break
			}
		}
		s.m.RUnlock()
		if idx == 0 { //  no expired subscribers
			continue
		}

		// remove expired subscribers
		s.m.Lock()
		log.Printf("drop %v subscribers", idx)
		for i := 0; i < idx; i++ {
			id := s.toDeletes[i].ID
			log.Printf("drop subscriber %v (%v)", id, s.subscribers[id].UserLogin)
		}

		A := s.toDeletes
		copy(A[0:len(A)-idx], A[idx:])
		s.toDeletes = s.toDeletes[:len(A)-idx]
		s.m.Unlock()
	}
}

func (s *EventStream) forward(event *Event) {
	s.m.RLock()
	defer s.m.RUnlock()

	log.Printf("try sending event %v", event.Type)
	for _, subscriber := range s.subscribers {
		log.Printf("subscriber %v (%v)", subscriber.ID, subscriber.UserLogin)
		if ShouldSendEvent(event, subscriber) {
			select {
			case subscriber.ch <- event:
				if event.Global {
					log.Printf("send event %v to all users", event.Type)
				} else {
					log.Printf("send event %v to %v", event.Type, event.UserLogin)
				}

			default:
				log.Println("out of channel buffer, drop event")
			}
		}
	}
}

func ShouldSendEvent(event *Event, subscriber *Subscriber) bool {
	return event.Global || event.UserLogin == subscriber.UserLogin
}

func (s *EventStream) SubscribeUser(id int64, userLogin string) (_ int64, ch chan *Event) {
	if id == 0 {
		panic("no id")
	}
	s.m.Lock()
	defer s.m.Unlock()
	if sub := s.subscribers[id]; sub != nil && sub.UserLogin == userLogin {
		log.Printf("reconnect to subscriber %v (%v)", id, sub.UserLogin)
		return id, sub.ch
	}
	log.Printf("create new subscriber %v (%v)", id, userLogin)
	sub := &Subscriber{
		ID:        id,
		AllEvents: true,
		UserLogin: userLogin,

		ch: make(chan *Event, 16),
	}
	s.subscribers[id] = sub
	return sub.ID, sub.ch
}

func (s *EventStream) Unsubscribe(id int64) {
	s.m.Lock()
	defer s.m.Unlock()

	t := time.Now().Add(time.Minute)
	A := s.toDeletes
	for i, d := range s.toDeletes {
		if d.ID == id {
			// move it to the back
			d.T = t
			copy(A[i:len(A)-1], A[i+1:])
			A[len(A)-1] = d
			return
		}
	}
	toDel := ToDelete{ID: id, T: t}
	s.toDeletes = append(s.toDeletes, toDel)
}

func (s *EventStream) Poll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, _ := session.GetBearer(r)
	ss, _ := s.Sessions.Authorize(ctx, token)
	if ss == nil {
		sutil.Respond(w).Error(401, errors.New("forbidden"))
		return
	}

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Println("invalid id", err)
		sutil.Respond(w).Error(400, errors.New("invalid id"))
		return
	}

	subscriberID, eventChannel := s.SubscribeUser(id, ss.Login)
	defer s.Unsubscribe(subscriberID)

	var events []interface{}
	var timer <-chan time.Time
	timeout := time.After(time.Minute)
	for {
		select {
		case event := <-eventChannel:
			log.Printf("subscriber %v (%v): prepare event %v `%v`", subscriberID, ss.Login, event.Type, event.Payload)
			events = append(events, event.Payload)
			if timer == nil {
				timer = time.After(100 * time.Millisecond)
			}

		case <-s.ctx.Done():
			log.Printf("server cancelled")
			return
		case <-ctx.Done():
			log.Printf("request cancelled")
			return
		case <-timer:
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(events)
			log.Printf("subscriber %v (%v): respond %v events", subscriberID, ss.Login, len(events))
			return
		case <-timeout:
			w.WriteHeader(200)
			_, _ = w.Write([]byte("[]"))
			log.Printf("subscriber %v (%v): respond timeout", subscriberID, ss.Login)
			return
		}
	}
}

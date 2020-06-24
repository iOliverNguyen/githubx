package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ng-vu/githubx/backend/cmd/githubx/config"
	"github.com/ng-vu/githubx/backend/pkg/eventstream"
	"github.com/ng-vu/githubx/backend/pkg/github"
	"github.com/ng-vu/githubx/backend/pkg/loader"
	"github.com/ng-vu/githubx/backend/pkg/session"
	"github.com/ng-vu/githubx/backend/pkg/store"
)

func main() {
	flConfigPath := flag.String("config-file", "", "path to config file")
	flag.Parse()

	cfg := config.Default()
	if *flConfigPath != "" {
		var err error
		cfg, err = config.Load(*flConfigPath)
		must(err)
	}
	if cfg.GitHub.APIKey == "" {
		panic("no apikey")
	}

	ctx := context.Background()
	logWriter := os.Stdout
	if cfg.LogFile != "" {
		f, err := os.Create(cfg.LogFile)
		must(err)
		logWriter = f
	}
	client := github.NewClient(cfg.OrgRepo, cfg.GitHub.APIKey, logWriter)
	{
		resp, err := client.Ping(ctx)
		must(err)
		log.Println("login as", resp.Viewer.Login)
	}

	st, err := store.New(store.Config{DBFile: cfg.DBFile})
	must(err)

	ld := loader.NewLoader(client, st, cfg.OrgRepo)
	go func() {
		must(ld.LoadAllIssues(ctx))
	}()

	sm := session.New(cfg.OrgRepo, client)
	es := eventstream.New(ctx, sm)
	s := NewService(cfg.OrgRepo, client, st, sm, es)

	m := http.NewServeMux()
	m.HandleFunc("/api/Authorize", s.ServeAuthorize)
	m.HandleFunc("/api/ListIssues", s.ServeListIssues)
	m.HandleFunc("/api/callback", s.ServeWebhook)
	m.HandleFunc("/api/poll", es.Poll)

	mux := http.NewServeMux()
	mux.Handle("/api/", CORS(m))
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

	log.Printf("http server listening at %v", cfg.Listen)
	if err = http.ListenAndServe(cfg.Listen, mux); err != http.ErrServerClosed {
		log.Printf("http server %v", err)
	}
}

func CORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Add("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

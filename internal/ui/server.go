package ui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nitschmann/hora/internal/service"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
	timeService service.TimeTracking
}

func NewServer(ts service.TimeTracking) *Server {
	return &Server{timeService: ts}
}

func (s *Server) Start(ctx context.Context, port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		limitStr := query.Get("limit")
		daysStr := query.Get("days")

		var limit int = 50
		var since *time.Time

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		if daysStr != "" {
			if days, err := strconv.Atoi(daysStr); err == nil {
				cutoff := time.Now().AddDate(0, 0, -days)
				since = &cutoff
			}
		}

		entries, err := s.timeService.GetAllEntriesWithPauses(ctx, limit, "desc", since)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			indexFile, err := staticFS.ReadFile("static/index.html")
			if err != nil {
				http.Error(w, "Index file not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexFile)

			return
		}

		// Serve other static files
		fs := http.FileServer(http.FS(staticFS))
		fs.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", port)

	fmt.Printf("WebUI running at http://%s\n", addr)

	return http.ListenAndServe(addr, mux)
}

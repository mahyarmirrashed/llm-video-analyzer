package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/cmd"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
)

type Server struct {
	cfg    *config.Config
	cmd    *cmd.Command
	Router *chi.Mux
}

func New(cfg *config.Config) (*Server, error) {
	db, err := qdrant.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	r := chi.NewRouter()
	s := &Server{
		cfg:    cfg,
		cmd:    cmd.New(cfg, db),
		Router: r,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

func (s *Server) setupMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) setupRoutes() {
	s.Router.Route("/api", func(r chi.Router) {
		r.Get("/health", s.handleHealth)
		r.Post("/process", s.handleProcess)
		r.Post("/search", s.handleSearch)
		r.Post("/clean", s.handleClean)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleProcess(w http.ResponseWriter, r *http.Request) {
	type processRequest struct {
		Url string `json:"url"`
	}

	var req processRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Url == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	_, err := s.cmd.Process(r.Context(), req.Url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process video")
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	type searchRequest struct {
		Query string `json:"query"`
		Limit int    `json:"limit,omitempty"`
	}

	var req searchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Query == "" {
		writeError(w, http.StatusBadRequest, "query is required")
		return
	}

	if req.Limit == 0 {
		req.Limit = s.cfg.QueryLimit
	}

	res, err := s.cmd.Query(r.Context(), req.Query, req.Limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	if err := s.cmd.Clean(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to clean database")
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

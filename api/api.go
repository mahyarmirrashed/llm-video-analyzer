package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/ollama"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/video"
)

type Server struct {
	cfg    *config.Config
	db     *qdrant.Client
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
		db:     db,
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
	err := r.ParseMultipartForm(300 << 20) // 300 MB
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		writeError(w, http.StatusBadRequest, "video file is required")
		return
	}
	defer file.Close()

	tempDir, err := os.MkdirTemp("", "llm-video-analzyer")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create temporary processing directory")
		return
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, header.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create temp file")
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save video file")
		return
	}

	v, err := video.New(tempFilePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize video")
		return
	}

	if err := v.Extract(s.cfg.SamplingInterval); err != nil {
		writeError(w, http.StatusInternalServerError, "frame extraction failed")
		return
	}
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

	desc, err := ollama.GetDescriptionFromQuery(r.Context(), s.cfg, req.Query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process query")
		return
	}

	embedding, err := ollama.GetTextEmbedding(r.Context(), s.cfg, desc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get embedding")
		return
	}

	res, err := s.db.Search(r.Context(), embedding, uint64(req.Limit))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Cleanup(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to clean database")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

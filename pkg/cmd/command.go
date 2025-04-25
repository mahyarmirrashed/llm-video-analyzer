package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/ollama"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/video"
)

type Command struct {
	cfg *config.Config
	db  *qdrant.Client
}

const downloadPath = "/tmp/llm-video-analyzer"

func New(cfg *config.Config, db *qdrant.Client) *Command {
	return &Command{
		cfg: cfg,
		db:  db,
	}
}

func (c *Command) Process(ctx context.Context, url string) (string, error) {
	d := video.NewYouTubeDownloader(downloadPath)

	path, err := d.Download(ctx, url)
	if err != nil {
		return "", err
	}
	defer os.Remove(path)

	v, err := video.New(path)
	if err != nil {
		return "", fmt.Errorf("failed to initialize video: %w", err)
	}

	if err := v.Extract(c.cfg.SamplingInterval); err != nil {
		return "", fmt.Errorf("frame extraction failed: %w", err)
	}
	defer v.Cleanup()

	for i := range v.Frames {
		frame := &v.Frames[i]

		if err := frame.Process(ctx, c.cfg); err != nil {
			log.Printf("skipping frame %s: %v", frame.Path, err)
			continue
		}

		if err := c.db.Store(ctx, url, frame); err != nil {
			log.Printf("failed to store frame %s: %v", frame.Path, err)
			continue
		}
	}

	return url, nil
}

func (c *Command) Query(ctx context.Context, query string, limit int) ([]qdrant.SearchResult, error) {
	desc, err := ollama.GetDescriptionFromQuery(ctx, c.cfg, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}

	embedding, err := ollama.GetTextEmbedding(ctx, c.cfg, desc)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}

	pts, err := c.db.Search(ctx, embedding, uint64(limit))
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	} else if len(pts) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return pts, nil
}

func (c *Command) Clean(ctx context.Context) error {
	err := c.db.Cleanup(ctx)
	if err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}

	return nil
}

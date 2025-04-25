package video

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/ollama"
)

type Frame struct {
	Path        string
	Timestamp   time.Duration
	Description string
	Embedding   []float32
}

func (f *Frame) Process(ctx context.Context, cfg *config.Config) error {
	log.Printf("processing frame at path: %s", f.Path)

	data, err := os.ReadFile(f.Path)
	if err != nil {
		return fmt.Errorf("failed to read frame: %w", err)
	}

	desc, err := ollama.GetDescriptionFromImage(ctx, cfg, data)
	if err != nil {
		return fmt.Errorf("failed to get description: %w", err)
	}
	f.Description = desc

	embedding, err := ollama.GetTextEmbedding(ctx, cfg, desc)
	if err != nil {
		return fmt.Errorf("failed to get embedding: %w", err)
	}
	f.Embedding = embedding

	log.Printf("processed frame, description: %q, embedding dim: %d", f.Description, len(f.Embedding))

	return nil
}

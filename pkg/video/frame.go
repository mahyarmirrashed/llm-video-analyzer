package video

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
)

type Frame struct {
	Path        string
	Timestamp   time.Duration
	Description string
	Embedding   []float32
}

func (f *Frame) Process(ctx context.Context, cfg *config.Config) error {
	log.Printf("processing frame at path: %s", f.Path)

	desc, err := f.getImageDescription(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to get description: %w", err)
	}
	f.Description = desc

	embedding, err := f.getTextEmbedding(ctx, cfg, desc)
	if err != nil {
		return fmt.Errorf("failed to get embedding: %w", err)
	}
	f.Embedding = embedding

	log.Printf("processed frame, description: %q, embedding dim: %d", f.Description, len(f.Embedding))

	return nil
}

func (f *Frame) getImageDescription(ctx context.Context, cfg *config.Config) (string, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read frame: %w", err)
	}

	payload := map[string]any{
		"model":  cfg.SamplingModel,
		"prompt": "Describe this video frame in detail for search purposes. Include objects, actions, colors, and context.",
		"stream": false,
		"images": []string{base64.StdEncoding.EncodeToString(data)},
	}

	rep, err := f.sendOllamaRequest(ctx, cfg, "/api/generate", payload)
	if err != nil {
		return "", err
	}

	var res struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(rep, &res); err != nil {
		return "", fmt.Errorf("failed to decode description response: %w", err)
	}

	return res.Response, nil
}

func (f *Frame) getTextEmbedding(ctx context.Context, cfg *config.Config, text string) ([]float32, error) {
	payload := map[string]any{
		"model":  cfg.EmbeddingModel,
		"prompt": text,
	}

	rep, err := f.sendOllamaRequest(ctx, cfg, "/api/embeddings", payload)
	if err != nil {
		return nil, err
	}

	var res struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(rep, &res); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	return res.Embedding, nil
}

func (f *Frame) sendOllamaRequest(ctx context.Context, cfg *config.Config, endpoint string, payload any) ([]byte, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s%s", cfg.OllamaURL, endpoint),
		bytes.NewBuffer(payloadJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}

	rep, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api request failed: %w", err)
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(rep.Body)
		return nil, fmt.Errorf("api error: %s (%d)", string(body), rep.StatusCode)
	}

	return io.ReadAll(rep.Body)
}

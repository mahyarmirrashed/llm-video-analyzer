package video

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

type ollamaResponse struct {
	Response  string    `json:"response"`
	Embedding []float32 `json:"embedding"`
}

func (f *Frame) Process(ctx context.Context, cfg *config.Config) error {
	image, err := os.ReadFile(f.Path)
	if err != nil {
		return fmt.Errorf("failed to read frame: %w", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(image)
	payload := map[string]interface{}{
		"model":  cfg.SamplingModel,
		"prompt": "Describe this video frame in detail for search purposes. Include objects, actions, colors, and context.",
		"images": []string{encodedImage},
		"options": map[string]interface{}{
			"embedding": true,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/generate", cfg.OllamaURL),
		bytes.NewBuffer(payloadJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}

	rep, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ollama api request failed: %w", err)
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(rep.Body)
		return fmt.Errorf("ollama api error: %s (%d)", string(body), rep.StatusCode)
	}

	var ollamaRep ollamaResponse
	if err := json.NewDecoder(rep.Body).Decode(&ollamaRep); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	f.Description = ollamaRep.Response
	f.Embedding = ollamaRep.Embedding

	return nil
}

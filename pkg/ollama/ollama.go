package ollama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
)

func GetImageDescription(ctx context.Context, cfg *config.Config, data []byte) (string, error) {
	payload := map[string]any{
		"model":  cfg.SamplingModel,
		"prompt": "Describe this video frame in detail for search purposes. Include objects, actions, colors, and context.",
		"stream": false,
		"images": []string{base64.StdEncoding.EncodeToString(data)},
	}

	rep, err := request(ctx, cfg, "/api/generate", payload)
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

func GetTextEmbedding(ctx context.Context, cfg *config.Config, text string) ([]float32, error) {
	payload := map[string]any{
		"model":  cfg.EmbeddingModel,
		"prompt": text,
	}

	rep, err := request(ctx, cfg, "/api/embeddings", payload)
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

func request(ctx context.Context, cfg *config.Config, endpoint string, payload any) ([]byte, error) {
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

	client := &http.Client{Timeout: 120 * time.Second}

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

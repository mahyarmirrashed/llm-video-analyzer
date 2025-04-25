package qdrant

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/video"
	"github.com/qdrant/go-client/qdrant"
)

type Client struct {
	*qdrant.Client
}

type SearchResult struct {
	Url         string
	Timestamp   float64
	Description string
	Score       float32
}

const (
	collectionName           = "llm-video-analyzer-frames"
	collectionDimensionality = 768
)

func New(databaseURL string) (*Client, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	host := u.Hostname()
	port := u.Port()

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: portInt,
	})
	if err != nil {
		return nil, err
	}

	res := Client{client}

	// ensure collection exists
	ctx := context.Background()
	exists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = res.createCollection(ctx, collectionName)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (c *Client) Cleanup(ctx context.Context) error {
	err := c.Client.DeleteCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	return c.createCollection(ctx, collectionName)
}

func (c *Client) Search(ctx context.Context, embedding []float32, limit uint64) ([]SearchResult, error) {
	if len(embedding) != collectionDimensionality {
		return nil, fmt.Errorf("embedding dimensions must be %d, got %d", collectionDimensionality, len(embedding))
	}
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	rep, err := c.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(embedding...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query points: %w", err)
	}

	res := make([]SearchResult, 0, len(rep))
	for _, pt := range rep {
		payload := pt.GetPayload()

		res = append(res, SearchResult{
			Url:         payload["url"].GetStringValue(),
			Timestamp:   payload["timestamp"].GetDoubleValue(),
			Description: payload["description"].GetStringValue(),
			Score:       pt.GetScore(),
		})
	}

	return res, nil
}

func (c *Client) Store(ctx context.Context, url string, frame *video.Frame) error {
	_, err := c.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{{
			Id:      qdrant.NewIDUUID(uuid.NewString()),
			Vectors: qdrant.NewVectors(frame.Embedding...),
			Payload: qdrant.NewValueMap(map[string]any{
				"url":         url,
				"timestamp":   frame.Timestamp.Seconds(),
				"description": frame.Description,
			}),
		}},
	})

	return err
}

func (c *Client) createCollection(ctx context.Context, collectionName string) error {
	err := c.Client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     collectionDimensionality,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return err
	}

	return nil
}

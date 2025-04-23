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

const collectionName = "llm-video-analyzer-frames"

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

	// ensure collection exists
	ctx := context.Background()
	exists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     7168,
				Distance: qdrant.Distance_Cosine,
			}),
		})

		if err != nil {
			return nil, err
		}
	}

	return &Client{client}, nil
}

func (c *Client) Store(ctx context.Context, videoID string, frame *video.Frame) error {
	_, err := c.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{{
			Id:      qdrant.NewIDUUID(uuid.NewString()),
			Vectors: qdrant.NewVectors(frame.Embedding...),
			Payload: qdrant.NewValueMap(map[string]interface{}{
				"video_id":    videoID,
				"timestamp":   frame.Timestamp.Seconds(),
				"description": frame.Description,
			}),
		}},
	})

	return err
}

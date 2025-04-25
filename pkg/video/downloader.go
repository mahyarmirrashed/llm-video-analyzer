package video

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type YouTubeDownloader struct {
	TempDir string
}

func NewYouTubeDownloader(tempDir string) *YouTubeDownloader {
	return &YouTubeDownloader{TempDir: tempDir}
}

func (yd *YouTubeDownloader) Download(ctx context.Context, url string) (string, error) {
	if err := os.MkdirAll(yd.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	tempFile := filepath.Join(yd.TempDir, "download.mp4")

	cmd := exec.CommandContext(ctx, "yt-dlp", "-o", tempFile, url)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	if _, err := os.Stat(tempFile); err != nil {
		return "", fmt.Errorf("downloaded file not found: %w", err)
	}

	return tempFile, nil
}

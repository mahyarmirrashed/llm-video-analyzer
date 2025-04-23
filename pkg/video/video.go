package video

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Video struct {
	ID             string
	Path           string
	Frames         []Frame
	ProcessingPath string
}

func New(path string) (*Video, error) {
	id, err := hash(path)
	if err != nil {
		return nil, fmt.Errorf("failed to hash video: %w", err)
	}

	return &Video{
		ID:   fmt.Sprintf("%x-%s", id, filepath.Base(path)),
		Path: path,
	}, nil
}

func (v *Video) Cleanup() error {
	if v.ProcessingPath == "" {
		return nil
	}

	return os.RemoveAll(v.ProcessingPath)
}

func (v *Video) Extract(interval int) error {
	v.ProcessingPath = filepath.Join(os.TempDir(), "llm-video-analyze", v.ID)
	if err := os.MkdirAll(v.ProcessingPath, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", v.Path,
		"-vf", fmt.Sprintf("fps=1/%d", interval),
		"-strftime", "1",
		filepath.Join(v.ProcessingPath, "frame_%Ts.png"), // %T is timestamp in seconds
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	frames, _ := filepath.Glob(filepath.Join(v.ProcessingPath, "frame_*.png"))
	v.Frames = make([]Frame, len(frames))
	for i, f := range frames {
		v.Frames[i] = Frame{
			Path:      f,
			Timestamp: parseTimestamp(f),
		}
	}

	return nil
}

func hash(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func parseTimestamp(path string) time.Duration {
	filename := filepath.Base(path)

	// Parses timestamps from FFmpeg's `-strftime` pattern
	if parts := strings.Split(filename, "_"); len(parts) > 1 {
		timePart := strings.TrimSuffix(parts[1], filepath.Ext(parts[1]))
		if secs, err := strconv.ParseFloat(strings.TrimSuffix(timePart, "s"), 64); err == nil {
			return time.Duration(secs * float64(time.Second))
		}
	}

	return 0
}

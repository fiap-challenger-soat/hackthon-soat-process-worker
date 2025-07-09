package processor_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/processor"
)

func TestFFmpegProcessor_Process(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tmpDir := t.TempDir()
		videoPath := filepath.Join(tmpDir, "test.mp4")

		cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=black:s=320x240:d=2", videoPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to create dummy video: %v, output: %s", err, output)
		}

		processor := processor.NewFFmpegProcessor()
		fullPath, fileName, err := processor.Process(context.Background(), videoPath)
		if err != nil {
			t.Fatalf("Process failed: %v", err)
		}
		if fullPath == "" || fileName == "" {
			t.Errorf("Expected non-empty fullPath and fileName, got '%s', '%s'", fullPath, fileName)
		}
		if !filepath.IsAbs(fullPath) {
			t.Errorf("Expected absolute path for zip, got: %s", fullPath)
		}
		if filepath.Ext(fileName) != ".zip" {
			t.Errorf("Expected .zip file, got: %s", fileName)
		}
	})

	t.Run("InvalidVideo", func(t *testing.T) {
		processor := processor.NewFFmpegProcessor()
		_, _, err := processor.Process(context.Background(), "nonexistent.mp4")
		if err == nil {
			t.Fatal("Expected error for nonexistent video file, got nil")
		}
		if !strings.Contains(err.Error(), "ffmpeg execution error") {
			t.Errorf("Expected ffmpeg execution error, got: %v", err)
		}
	})

	t.Run("ContextCancel", func(t *testing.T) {
		tmpDir := t.TempDir()
		videoPath := filepath.Join(tmpDir, "test.mp4")
		cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=black:s=320x240:d=2", videoPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to create dummy video: %v, output: %s", err, output)
		}

		processor := processor.NewFFmpegProcessor()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := processor.Process(ctx, videoPath)
		if err == nil {
			t.Fatal("Expected error due to context cancellation, got nil")
		}
	})
}

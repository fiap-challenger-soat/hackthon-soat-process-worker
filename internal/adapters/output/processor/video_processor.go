package processor

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type ffmpegProcessor struct{}

func NewFFmpegProcessor() *ffmpegProcessor {
	return &ffmpegProcessor{}
}

func (p *ffmpegProcessor) Process(ctx context.Context, localVideoPath string) (string, string, error) {
	frameDir, err := os.MkdirTemp("", "frames-*")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp dir for frames: %w", err)
	}

	framePattern := filepath.Join(frameDir, "frame_%04d.png")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localVideoPath, "-vf", "fps=1", framePattern)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("ffmpeg execution error: %w - output: %s", err, string(output))
	}

	return zipFrames(frameDir)
}

func zipFrames(sourceDir string) (string, string, error) {
	zipFile, err := os.CreateTemp("", "archive-*.zip")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp zip file: %w", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	frames, err := filepath.Glob(filepath.Join(sourceDir, "*.png"))
	if err != nil {
		return "", "", fmt.Errorf("failed to find frames: %w", err)
	}
	if len(frames) == 0 {
		return "", "", fmt.Errorf("no frames extracted")
	}

	for _, framePath := range frames {
		if err := addFileToZip(zipWriter, framePath); err != nil {
			return "", "", fmt.Errorf("failed to add '%s' to zip: %w", framePath, err)
		}
	}

	fullPath := zipFile.Name()
	fileName := filepath.Base(fullPath)
	return fullPath, fileName, nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

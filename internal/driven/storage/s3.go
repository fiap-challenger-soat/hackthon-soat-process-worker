package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/model"
)

type S3FileHandler interface {
	DownloadFile(ctx context.Context, objectKey string) (*model.DownloadedFile, error)
	UploadFile(ctx context.Context, localFilePath, objectKey string) error
}

type s3client struct {
	client         *s3.Client
	bucketDownName string
	bucketUpName   string
}

func NewS3Adapter(s3Client *s3.Client, bucketUpName, bucketDownName string) *s3client {
	return &s3client{
		client:         s3Client,
		bucketUpName:   bucketDownName,
		bucketDownName: bucketUpName,
	}
}

func (a *s3client) DownloadFile(ctx context.Context, objectKey string) (*model.DownloadedFile, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.bucketDownName),
		Key:    aws.String(objectKey),
	}

	result, err := a.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object '%s' from S3: %w", objectKey, err)
	}

	tempFile, err := os.CreateTemp("", "video-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := io.Copy(tempFile, result.Body); err != nil {
		return nil, fmt.Errorf("failed to copy S3 content: %w", err)
	}

	downloaded := &model.DownloadedFile{
		Path: tempFile.Name(),
		File: tempFile,
	}

	return downloaded, nil
}

func (a *s3client) UploadFile(ctx context.Context, localFilePath, objectKey string) error {

	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s' for upload: %w", localFilePath, err)
	}

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for '%s': %w", localFilePath, err)
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(a.bucketUpName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(stat.Size()),
		ContentType:   aws.String("application/octet-stream"),
	}

	_, err = a.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload object '%s' to S3: %w", objectKey, err)
	}

	return nil
}

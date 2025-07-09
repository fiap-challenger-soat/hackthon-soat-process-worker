package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/input"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/processor"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/queue"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/repository"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/storage"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/clients/aws"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/clients/postgres"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/config"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/service"
)

func main() {
	log.Println("INFO: Starting the worker service...")

	// Load configuration
	config.Init()
	cfg := config.Vars

	ctx := context.Background()

	// Initialize clients
	db, err := postgres.NewPostgresClient()
	if err != nil {
		log.Fatalf("FATAL: erro ao conectar ao banco: %v", err)
	}

	awsCfg, err := aws.NewAWSConfig(ctx)
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to initialize AWS configuration: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = true })
	sqsClient := sqs.NewFromConfig(awsCfg)

	// Initialize adapters
	videoRepository := repository.NewVideoJobRepository(db)
	storageAdapter := storage.NewS3Adapter(s3Client, cfg.S3UploadBucket, cfg.S3DownloadBucket)
	videoProcessingAdapter := processor.NewFFmpegProcessor()
	sqsMessageQueueAdapter := queue.NewSQSAdapter(sqsClient, cfg.SQSErrorQueueURL, cfg.SQSWorkQueueURL)

	// Initialize service and consumer
	jobService := service.NewJobService(
		videoRepository,
		storageAdapter,
		videoProcessingAdapter,
		sqsMessageQueueAdapter,
	)

	sqsConsumer := input.NewConsumer(sqsMessageQueueAdapter, jobService)

	// Start consumer
	go sqsConsumer.Start(ctx)
	select {}
}

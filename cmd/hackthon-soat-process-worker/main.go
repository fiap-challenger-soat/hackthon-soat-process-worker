package main

import (
	"context"
	"log"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/config"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/clients/aws"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/clients/postgres"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/clients/redis"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/service"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/cache"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/processor"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/queue"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/repository"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/storage"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driver/consumer"
)

func main() {
	log.Println("INFO: Starting the worker service...")

	// Load configuration
	config.Init()
	cfg := config.Vars

	ctx := context.Background()

	// Initialize clients
	dbClient, err := postgres.NewPostgresClient()
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to initialize Postgres client: %v", err)
	}

	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to initialize Redis client: %v", err)
	}

	awsCfg, err := aws.NewAwsConfig(ctx)
	if err != nil {
		log.Fatalf("FATAL ERROR: Failed to initialize AWS configuration: %v", err)
	}
	s3Client := aws.NewS3Client(awsCfg)
	sqsClient := aws.NewSQSClient(awsCfg)

	// Initialize adapters
	videoRepository := repository.NewVideoJobRepository(dbClient)
	redisAdapter := cache.NewRedisAdapter(redisClient)
	storageAdapter := storage.NewS3Adapter(s3Client, cfg.S3UploadBucket, cfg.S3DownloadBucket)
	videoProcessingAdapter := processor.NewFFmpegProcessor()
	sqsMessageQueueAdapter := queue.NewSQSAdapter(sqsClient, cfg.SQSErrorQueueURL, cfg.SQSWorkQueueURL)

	// Initialize service and consumer
	jobService := service.NewJobService(
		videoRepository,
		redisAdapter,
		storageAdapter,
		videoProcessingAdapter,
		sqsMessageQueueAdapter,
	)
	sqsConsumer := consumer.NewConsumer(sqsMessageQueueAdapter, jobService)

	// Start consumer
	go sqsConsumer.Start(ctx)
	select {}
}

package ports

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
)

//go:generate mockgen -destination=mocks/mock_s3client.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports S3Client
type S3Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

//go:generate mockgen -destination=mocks/mock_sqsclient.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports SQSClient
type SQSClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

//go:generate mockgen -destination=mocks/mock_videojobrepository.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports VideoJobRepository
type VideoJobRepository interface {
	GetJobByID(ctx context.Context, jobID string) (*domain.VideoJobDTO, error)
	UpdateJobStatus(ctx context.Context, videoJob *domain.VideoJob) error
}

//go:generate mockgen -destination=mocks/mock_s3adapter.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports S3Adapter
type S3Adapter interface {
	DownloadFile(ctx context.Context, objectKey string) (*domain.DownloadedFile, error)
	UploadFile(ctx context.Context, localFilePath, objectKey string) error
}

//go:generate mockgen -destination=mocks/mock_sqsadapter.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports SQSAdapter
type SQSAdapter interface {
	Publish(ctx context.Context, event domain.JobErrorEvent) error
	Receive(ctx context.Context, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error)
	Delete(ctx context.Context, receiptHandle string) error
}

//go:generate mockgen -destination=mocks/mock_processoradapter.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports ProcessorAdapter
type ProcessorAdapter interface {
	Process(ctx context.Context, localVideoPath string) (string, string, error)
}

//go:generate mockgen -destination=mocks/mock_jobservice.go -package=mocks github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports JobService
type JobService interface {
	ProcessJob(ctx context.Context, jobID, videoPath string) error
}

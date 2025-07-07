package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/model"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/cache"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/processor"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/queue"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/repository"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/driven/storage"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type JobProcessor interface {
	ProcessJob(ctx context.Context, jobID, videoPath string) error
}

type JobService struct {
	repo      repository.VideoJobRepository
	cache     cache.JobCache
	storage   storage.S3FileHandler
	processor processor.WorkProcessor
	errorPub  queue.WorkQueue
}

func NewJobService(
	repo repository.VideoJobRepository,
	cache cache.JobCache,
	storage storage.S3FileHandler,
	processor processor.WorkProcessor,
	errorPub queue.WorkQueue,
) *JobService {
	return &JobService{
		repo:      repo,
		cache:     cache,
		storage:   storage,
		processor: processor,
		errorPub:  errorPub,
	}
}

func (s *JobService) ProcessJob(ctx context.Context, jobID, videoPath string) error {
	log.Printf("[Job %s] Starting processing for video: %s", jobID, videoPath)

	if videoPath == "" {
		log.Printf("CRITICAL ERROR: [Job %s] videoPath is empty. Invalid message will be discarded.", jobID)
		return nil
	}

	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: [Job %s] Job not found in DB. Message discarded.", jobID)
			return nil
		}
		return fmt.Errorf("job %s: failed to fetch job details: %w", jobID, err)
	}

	if err := s.setStatus(ctx, job, model.VideoStatusProcessing); err != nil {
		return fmt.Errorf("job %s: failed to update status to 'processing': %w", jobID, err)
	}

	log.Printf("[Job %s] Downloading video from S3...", jobID)
	tempVideoFile, err := s.storage.DownloadFile(ctx, videoPath)
	if err != nil {
		return fmt.Errorf("job %s: failed to download video from S3: %w", jobID, err)
	}

	log.Printf("[Job %s] Processing local temporary file: %s", jobID, tempVideoFile.Path)
	localZipPath, zipName, err := s.processor.Process(ctx, tempVideoFile.Path)
	if err != nil {
		s.handleProcessingError(ctx, job)
		return nil
	}

	outputPath := fmt.Sprintf("output/%s", zipName)
	log.Printf("[Job %s] Uploading compressed file to S3: %s", jobID, outputPath)
	if err := s.storage.UploadFile(ctx, localZipPath, outputPath); err != nil {
		s.handleProcessingError(ctx, job)
		return nil
	}

	job.OutputPath = lo.ToPtr(outputPath)
	if err := s.setStatus(ctx, job, model.VideoStatusCompleted); err != nil {
		return fmt.Errorf("job %s: job completed, but failed to update final status: %w", job.ID, err)
	}

	log.Printf("[Job %s] Processing completed successfully.", jobID)
	return nil
}

func (s *JobService) handleProcessingError(ctx context.Context, job *model.VideoJob) {
	log.Printf("[Job %s] ERROR: failed to process video", job.ID)
	_ = s.setStatus(ctx, job, model.VideoStatusFailed)

	errorEvent := model.JobErrorEvent{
		JobID:    job.ID,
		UserID:   job.UserID,
		FailedAt: time.Now().Format(time.RFC3339),
	}
	if err := s.errorPub.Publish(ctx, errorEvent); err != nil {
		log.Printf("CRITICAL ERROR: [Job %s] Failed to publish to error queue: %v", job.ID, err)
	}
}

func (s *JobService) setStatus(ctx context.Context, job *model.VideoJob, status model.VideoStatus) error {
	job.Status = status
	history := model.JobStatusHistory{
		ID:     uuid.New().String(),
		JobID:  job.ID,
		Status: status,
	}
	if err := s.repo.UpdateJobStatus(ctx, job, &history); err != nil {
		return err
	}
	if err := s.cache.SetJobStatus(job.ID, string(status)); err != nil {
		log.Printf("WARNING: [Job %s] Failed to update Redis cache: %v", job.ID, err)
	}
	return nil
}

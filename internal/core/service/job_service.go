package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type JobService struct {
	repo      ports.VideoJobRepository
	storage   ports.S3Adapter
	processor ports.ProcessorAdapter
	errorPub  ports.SQSAdapter
}

func NewJobService(
	repo ports.VideoJobRepository,
	storage ports.S3Adapter,
	processor ports.ProcessorAdapter,
	errorPub ports.SQSAdapter,
) *JobService {
	return &JobService{
		repo:      repo,
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

	if err := s.setStatus(ctx, job, domain.VideoStatusProcessing); err != nil {
		return fmt.Errorf("job %s: failed to update status to 'processing': %w", jobID, err)
	}

	tempVideoFile, err := s.storage.DownloadFile(ctx, videoPath)
	if err != nil {
		s.failJob(ctx, job)
		return fmt.Errorf("job %s: failed to download video from S3: %w", jobID, err)
	}

	localZipPath, zipName, err := s.processor.Process(ctx, tempVideoFile.Path)
	if err != nil {
		s.failJob(ctx, job)
		return nil
	}

	outputPath := fmt.Sprintf("output/%s", zipName)
	if err := s.storage.UploadFile(ctx, localZipPath, outputPath); err != nil {
		s.failJob(ctx, job)
		return nil
	}

	job.OutputPath = lo.ToPtr(outputPath)
	if err := s.setStatus(ctx, job, domain.VideoStatusCompleted); err != nil {
		return fmt.Errorf("job %s: job completed, but failed to update final status: %w", job.ID, err)
	}

	log.Printf("[Job %s] Processing completed successfully.", jobID)
	return nil
}

func (s *JobService) failJob(ctx context.Context, job *domain.VideoJobDTO) {
	log.Printf("[Job %s] ERROR: failed to process video", job.ID)
	_ = s.setStatus(ctx, job, domain.VideoStatusFailed)
	event := domain.JobErrorEvent{
		JobID: job.ID,
		Email: job.Email,
	}
	if err := s.errorPub.Publish(ctx, event); err != nil {
		log.Printf("CRITICAL ERROR: [Job %s] Failed to publish to error queue: %v", job.ID, err)
	}
}

func (s *JobService) setStatus(ctx context.Context, job *domain.VideoJobDTO, status domain.VideoStatus) error {
	return s.repo.UpdateJobStatus(ctx, &domain.VideoJob{
		ID:         job.ID,
		Status:     status,
		CreatedAt:  job.CreatedAt,
		OutputPath: job.OutputPath,
		UserID:     job.UserID,
		VideoPath:  job.VideoPath,
	})
}

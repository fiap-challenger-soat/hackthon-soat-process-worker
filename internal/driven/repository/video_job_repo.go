package repository

import (
	"context"
	"fmt"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/model"
	"gorm.io/gorm"
)

type VideoJobRepository interface {
	GetJobByID(ctx context.Context, jobID string) (*model.VideoJob, error)
	UpdateJobStatus(ctx context.Context, videoJob *model.VideoJob, historyJob *model.JobStatusHistory) error
}

type videoJobRepository struct {
	db *gorm.DB
}

func NewVideoJobRepository(db *gorm.DB) VideoJobRepository {
	return &videoJobRepository{
		db: db,
	}
}

func (r *videoJobRepository) GetJobByID(ctx context.Context, jobID string) (*model.VideoJob, error) {
	var job model.VideoJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job with id '%s' not found", jobID)
		}
		return nil, fmt.Errorf("error fetching job with id '%s': %w", jobID, err)
	}
	return &job, nil
}

func (r *videoJobRepository) UpdateJobStatus(ctx context.Context, videoJob *model.VideoJob, historyJob *model.JobStatusHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(videoJob).Error; err != nil {
			return fmt.Errorf("failed to update video job: %w", err)
		}
		if err := tx.Create(historyJob).Error; err != nil {
			return fmt.Errorf("failed to create job status history: %w", err)
		}
		return nil
	})
}

package repository

import (
	"context"
	"fmt"

	model "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"gorm.io/gorm"
)

type videoJobRepository struct {
	db *gorm.DB
}

func NewVideoJobRepository(db *gorm.DB) *videoJobRepository {
	return &videoJobRepository{db: db}
}
func (r *videoJobRepository) GetJobByID(ctx context.Context, jobID string) (*model.VideoJobDTO, error) {
	var job model.VideoJobDTO
	err := r.db.WithContext(ctx).
		Table("tb_video_jobs").
		Select("tb_video_jobs.id, tb_video_jobs.status, tb_video_jobs.created_at, tb_video_jobs.output_path, tb_video_jobs.user_id, tb_video_jobs.video_path, tb_user.email").
		Joins("join tb_user on tb_user.id = tb_video_jobs.user_id").
		Where("tb_video_jobs.id = ?", jobID).
		First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job with id '%s' not found", jobID)
		}
		return nil, fmt.Errorf("error fetching job with id '%s': %w", jobID, err)
	}
	return &job, nil
}

func (r *videoJobRepository) UpdateJobStatus(ctx context.Context, videoJob *model.VideoJob) error {
	if err := r.db.WithContext(ctx).Save(videoJob).Error; err != nil {
		return fmt.Errorf("failed to update video job: %w", err)
	}
	return nil
}

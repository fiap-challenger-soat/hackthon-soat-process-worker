package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/repository"
	model "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/pkg/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type repositoryTestSuite struct {
	suite.Suite
	mockDB   *gorm.DB
	mockSQL  sqlmock.Sqlmock
	repo     ports.VideoJobRepository
	ctx      context.Context
	videoDTO model.VideoJobDTO
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(repositoryTestSuite))
}

func (rts *repositoryTestSuite) SetupTest() {
	rts.mockSQL, rts.mockDB = tests.BuildMockDB(rts.T())
}

func (rts *repositoryTestSuite) AfterTest(_, _ string) {
	assert.NoError(rts.T(), rts.mockSQL.ExpectationsWereMet())
}

func (rts *repositoryTestSuite) BeforeTest(suiteName, testName string) {
	rts.ctx = context.Background()
	var err error
	rts.repo = repository.NewVideoJobRepository(rts.mockDB)
	assert.NoError(rts.T(), err)
	rts.videoDTO = model.VideoJobDTO{
		ID:         "194f2506-3a19-42fb-91a0-50442a1bfcfd",
		Status:     "processing",
		CreatedAt:  time.Now().String(),
		OutputPath: lo.ToPtr("donwloads/output.mp4"),
		UserID:     "51b00f47-5293-4100-8a24-0c5349b5ff47",
		VideoPath:  "uploads/video.mp4",
		Email:      "pedro@test.com",
	}
}

func (rts *repositoryTestSuite) Test_GetJobByID() {
	rts.T().Run("Should get a job by ID", func(t *testing.T) {
		const sqlRegexp = `(?i)SELECT .*FROM .*tb_video_jobs.*join.*tb_user.*WHERE.*tb_video_jobs.id.*`

		rows := sqlmock.NewRows([]string{
			"id", "status", "created_at", "output_path", "user_id", "video_path", "email",
		}).AddRow(
			rts.videoDTO.ID,
			rts.videoDTO.Status,
			rts.videoDTO.CreatedAt,
			rts.videoDTO.OutputPath,
			rts.videoDTO.UserID,
			rts.videoDTO.VideoPath,
			rts.videoDTO.Email,
		)
		rts.mockSQL.ExpectQuery(sqlRegexp).
			WithArgs(rts.videoDTO.ID, 1).
			WillReturnRows(rows)

		job, err := rts.repo.GetJobByID(rts.ctx, rts.videoDTO.ID)
		assert.NoError(t, err)
		assert.NotNil(t, job)
		assert.Equal(t, rts.videoDTO.ID, job.ID)

	})

	rts.T().Run("Should return not found error when job does not exist", func(t *testing.T) {
		const sqlRegexp = `(?i)SELECT .*FROM .*tb_video_jobs.*join.*tb_user.*WHERE.*tb_video_jobs.id.*`
		rts.mockSQL.ExpectQuery(sqlRegexp).
			WithArgs(rts.videoDTO.ID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		job, err := rts.repo.GetJobByID(rts.ctx, rts.videoDTO.ID)
		assert.Nil(t, job)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	rts.T().Run("Should return error when db returns unexpected error", func(t *testing.T) {
		const sqlRegexp = `(?i)SELECT .*FROM .*tb_video_jobs.*join.*tb_user.*WHERE.*tb_video_jobs.id.*`
		dbErr := fmt.Errorf("db connection error")
		rts.mockSQL.ExpectQuery(sqlRegexp).
			WithArgs(rts.videoDTO.ID, 1).
			WillReturnError(dbErr)

		job, err := rts.repo.GetJobByID(rts.ctx, rts.videoDTO.ID)
		assert.Nil(t, job)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching job")
		assert.ErrorIs(t, err, dbErr)
	})
}

func (rts *repositoryTestSuite) Test_UpdateJobStatus() {
	rts.T().Run("Should update job status successfully", func(t *testing.T) {
		videoJob := &model.VideoJob{
			ID:        rts.videoDTO.ID,
			Status:    "completed",
			CreatedAt: time.Now().String(),
			UserID:    rts.videoDTO.UserID,
			VideoPath: rts.videoDTO.VideoPath,
		}

		const sqlRegexp = `(?i)UPDATE .*tb_video_jobs.*SET.*status.*WHERE.*id.*`
		rts.mockSQL.ExpectBegin()
		rts.mockSQL.ExpectExec(sqlRegexp).
			WithArgs(
				videoJob.Status,
				videoJob.CreatedAt,
				videoJob.OutputPath,
				videoJob.UserID,
				videoJob.VideoPath,
				videoJob.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		rts.mockSQL.ExpectCommit()

		err := rts.repo.UpdateJobStatus(rts.ctx, videoJob)
		assert.NoError(t, err)
	})

	rts.T().Run("Should return error when db returns error", func(t *testing.T) {
		videoJob := &model.VideoJob{
			ID:        rts.videoDTO.ID,
			Status:    "failed",
			CreatedAt: time.Now().String(),
			UserID:    rts.videoDTO.UserID,
			VideoPath: rts.videoDTO.VideoPath,
		}

		const sqlRegexp = `(?i)UPDATE .*tb_video_jobs.*SET.*status.*WHERE.*id.*`
		dbErr := fmt.Errorf("db update error")
		rts.mockSQL.ExpectBegin()
		rts.mockSQL.ExpectExec(sqlRegexp).
			WithArgs(
				videoJob.Status,
				videoJob.CreatedAt,
				videoJob.OutputPath,
				videoJob.UserID,
				videoJob.VideoPath,
				videoJob.ID,
			).
			WillReturnError(dbErr)
		rts.mockSQL.ExpectRollback()

		err := rts.repo.UpdateJobStatus(rts.ctx, videoJob)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update video job")
		assert.ErrorIs(t, err, dbErr)
	})
}

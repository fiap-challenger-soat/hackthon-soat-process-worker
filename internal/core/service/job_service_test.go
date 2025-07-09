package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports/mocks"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/service"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type jobServiceTestSuite struct {
	suite.Suite

	ctx           context.Context
	mockRepo      *mocks.MockVideoJobRepository
	mockStorage   *mocks.MockS3Adapter
	mockProcessor *mocks.MockProcessorAdapter
	mockErrorPub  *mocks.MockSQSAdapter
	jobService    *service.JobService
}

func (sts *jobServiceTestSuite) BeforeTest(_, _ string) {
	ctrl := gomock.NewController(sts.T())
	sts.ctx = context.Background()
	sts.mockRepo = mocks.NewMockVideoJobRepository(ctrl)
	sts.mockStorage = mocks.NewMockS3Adapter(ctrl)
	sts.mockProcessor = mocks.NewMockProcessorAdapter(ctrl)
	sts.mockErrorPub = mocks.NewMockSQSAdapter(ctrl)
	sts.jobService = service.NewJobService(
		sts.mockRepo,
		sts.mockStorage,
		sts.mockProcessor,
		sts.mockErrorPub,
	)
}

func Test_JobServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(jobServiceTestSuite))
}

func (sts *jobServiceTestSuite) Test_ProcessJob_Sucess() {
	s := sts.T()

	s.Run("should process job successfully", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "started",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
		}

		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, j *domain.VideoJob) error {
			j.Status = domain.VideoStatusProcessing
			return nil
		})

		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(&domain.DownloadedFile{
			Path: "/testdata/downloadFile/trailerGTA6_4k.mp4",
			File: nil, // Assuming the file is not needed for this test
		}, nil)

		sts.mockProcessor.EXPECT().Process(sts.ctx, "/testdata/downloadFile/trailerGTA6_4k.mp4").Return("/testdata/processed/trailerGTA6_4k.zip", "trailerGTA6_4k.zip", nil)
		sts.mockStorage.EXPECT().UploadFile(sts.ctx, "/testdata/processed/trailerGTA6_4k.zip", "output/trailerGTA6_4k.zip").Return(nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, j *domain.VideoJob) error {
			j.Status = domain.VideoStatusCompleted
			j.OutputPath = lo.ToPtr("s3://donwload/trailerGTA6_4k.zip")
			return nil
		})

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error when processing job")

	})

}

func (sts *jobServiceTestSuite) Test_ProcessJob_Errors() {
	s := sts.T()

	s.Run("should return nil if videoPath is empty", func(t *testing.T) {
		jobID := "job-123"
		videoPath := ""

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error when videoPath is empty")
	})

	s.Run("should return nil if job not found in DB", func(t *testing.T) {
		jobID := "job-404"
		videoPath := "s3://upload/video.mp4"

		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(nil, gorm.ErrRecordNotFound)

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error when job is not found")
	})

	s.Run("should return error if GetJobByID fails with unexpected error", func(t *testing.T) {
		jobID := "job-err"
		videoPath := "s3://upload/video.mp4"
		expectedErr := errors.New("db error")

		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(nil, expectedErr)

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.Error(err)
		sts.Contains(err.Error(), "failed to fetch job details")
	})

	s.Run("should return error if setStatus to processing fails", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "started",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
		}
		expectedErr := errors.New("update status error")

		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(expectedErr)

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.Error(err)
		sts.Contains(err.Error(), "failed to update status to 'processing'")
	})

	s.Run("should return error if DownloadFile fails", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "started",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
		}
		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(nil, errors.New("download error"))

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.Error(err)
		sts.Contains(err.Error(), "failed to download video from S3")
	})

	s.Run("should fail job and return nil if processor fails", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
			Email:     "user@email.com",
		}
		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(&domain.DownloadedFile{
			Path: "/tmp/video.mp4",
			File: nil,
		}, nil)
		sts.mockProcessor.EXPECT().Process(sts.ctx, "/tmp/video.mp4").Return("", "", errors.New("process error"))
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockErrorPub.EXPECT().Publish(sts.ctx, gomock.Any()).Return(nil)

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error when processor fails (job is failed internally)")
	})

	s.Run("should fail job and return nil if upload fails", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
			Email:     "user@email.com",
		}
		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(&domain.DownloadedFile{
			Path: "/tmp/video.mp4",
			File: nil,
		}, nil)
		sts.mockProcessor.EXPECT().Process(sts.ctx, "/tmp/video.mp4").Return("/tmp/video.zip", "video.zip", nil)
		sts.mockStorage.EXPECT().UploadFile(sts.ctx, "/tmp/video.zip", "output/video.zip").Return(errors.New("upload error"))
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockErrorPub.EXPECT().Publish(sts.ctx, gomock.Any()).Return(nil)

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error when upload fails (job is failed internally)")
	})

	s.Run("should return error if setStatus to completed fails", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
		}
		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(&domain.DownloadedFile{
			Path: "/tmp/video.mp4",
			File: nil,
		}, nil)
		sts.mockProcessor.EXPECT().Process(sts.ctx, "/tmp/video.mp4").Return("/tmp/video.zip", "video.zip", nil)
		sts.mockStorage.EXPECT().UploadFile(sts.ctx, "/tmp/video.zip", "output/video.zip").Return(nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(errors.New("final status error"))

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.Error(err)
		sts.Contains(err.Error(), "job completed, but failed to update final status")
	})

	s.Run("should handle error publishing to error queue", func(t *testing.T) {
		jobID := "job-123"
		videoPath := "s3://upload/video.mp4"
		job := &domain.VideoJobDTO{
			ID:        jobID,
			Status:    "started",
			CreatedAt: "2023-10-01T00:00:00Z",
			UserID:    "user-123",
			VideoPath: videoPath,
			Email:     "pedrinho@gmail.com",
		}

		event := domain.JobErrorEvent{
			JobID: jobID,
			Email: "pedrinho@gmail.com",
		}

		sts.mockRepo.EXPECT().GetJobByID(sts.ctx, jobID).Return(job, nil)
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockStorage.EXPECT().DownloadFile(sts.ctx, videoPath).Return(&domain.DownloadedFile{
			Path: "/tmp/video.mp4",
			File: nil,
		}, nil)
		sts.mockProcessor.EXPECT().Process(sts.ctx, "/tmp/video.mp4").Return("/tmp/video.zip", "video.zip", nil)
		sts.mockStorage.EXPECT().UploadFile(sts.ctx, "/tmp/video.zip", "output/video.zip").Return(errors.New("upload error"))
		sts.mockRepo.EXPECT().UpdateJobStatus(sts.ctx, gomock.Any()).Return(nil)
		sts.mockErrorPub.EXPECT().Publish(sts.ctx, event).Return(errors.New("publish error"))

		err := sts.jobService.ProcessJob(sts.ctx, jobID, videoPath)

		sts.NoError(err, "expected no error even if publishing to error queue fails")
	})

}

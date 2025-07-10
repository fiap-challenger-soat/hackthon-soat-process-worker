package input_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/input"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type consumerTestSuite struct {
	suite.Suite

	ctx         context.Context
	mockQueue   *mocks.MockSQSAdapter
	mockService *mocks.MockJobService
	consumer    *input.Consumer
}

func (suite *consumerTestSuite) BeforeTest(_, _ string) {
	ctrl := gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.mockQueue = mocks.NewMockSQSAdapter(ctrl)
	suite.mockService = mocks.NewMockJobService(ctrl)
	suite.consumer = input.NewConsumer(suite.mockQueue, suite.mockService)
}

func Test_ConsumerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(consumerTestSuite))
}

func (suite *consumerTestSuite) Test_Start_Success() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Mock message to be received
	mockMsg := struct {
		jobID     string
		videoPath string
		body      string
		receipt   string
	}{
		jobID:     "job-123",
		videoPath: "upload/video.mp4",
		body:      `{"job_id":"job-123","video_path":"upload/video.mp4"}`,
		receipt:   "receipt-handle-123",
	}

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{
			{
				Body:          &mockMsg.body,
				ReceiptHandle: &mockMsg.receipt,
			},
		}, nil).
		Times(1)

	
	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	suite.mockService.EXPECT().
		ProcessJob(gomock.Any(), mockMsg.jobID, mockMsg.videoPath).
		Return(nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Delete(gomock.Any(), mockMsg.receipt).
		Return(nil).
		Times(1)

	
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

func (suite *consumerTestSuite) Test_Start_ReceiveError() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	
	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return(nil, assert.AnError).
		Times(1)

	
	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

func (suite *consumerTestSuite) Test_Start_InvalidMessage_NilBodyOrReceipt() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgs := []types.Message{
		{Body: nil, ReceiptHandle: nil},
		{Body: nil, ReceiptHandle: lo.ToPtr("receip")},
		{Body: lo.ToPtr(`{"job_id":"job-1","video_path":"foo.mp4"}`), ReceiptHandle: nil},
	}

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return(msgs, nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

func (suite *consumerTestSuite) Test_Start_InvalidJSONBody() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body := "not-a-json"
	receipt := "receipt-err"
	msgs := []types.Message{
		{Body: &body, ReceiptHandle: &receipt},
	}

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return(msgs, nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	suite.mockQueue.EXPECT().
		Delete(gomock.Any(), receipt).
		Return(nil).
		Times(1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

func (suite *consumerTestSuite) Test_Start_ProcessJobError() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockMsg := struct {
		jobID     string
		videoPath string
		body      string
		receipt   string
	}{
		jobID:     "job-err",
		videoPath: "video-err.mp4",
		body:      `{"job_id":"job-err","video_path":"video-err.mp4"}`,
		receipt:   "receipt-handle-err",
	}

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{
			{
				Body:          &mockMsg.body,
				ReceiptHandle: &mockMsg.receipt,
			},
		}, nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	suite.mockService.EXPECT().
		ProcessJob(gomock.Any(), mockMsg.jobID, mockMsg.videoPath).
		Return(assert.AnError).
		Times(1)

	suite.mockQueue.EXPECT().
		Delete(gomock.Any(), mockMsg.receipt).
		Return(nil).
		Times(1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

func (suite *consumerTestSuite) Test_Start_QueueDeleteError() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockMsg := struct {
		jobID     string
		videoPath string
		body      string
		receipt   string
	}{
		jobID:     "job-del-err",
		videoPath: "video-del-err.mp4",
		body:      `{"job_id":"job-del-err","video_path":"video-del-err.mp4"}`,
		receipt:   "receipt-handle-del-err",
	}

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{
			{
				Body:          &mockMsg.body,
				ReceiptHandle: &mockMsg.receipt,
			},
		}, nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Receive(gomock.Any(), int32(10), int32(20)).
		Return([]types.Message{}, nil).
		AnyTimes()

	suite.mockService.EXPECT().
		ProcessJob(gomock.Any(), mockMsg.jobID, mockMsg.videoPath).
		Return(nil).
		Times(1)

	suite.mockQueue.EXPECT().
		Delete(gomock.Any(), mockMsg.receipt).
		Return(assert.AnError).
		Times(1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	suite.consumer.Start(ctx)
}

package queue_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/adapters/output/queue"
	model "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type sqsHandleTest struct {
	suite.Suite

	ctx           context.Context
	sqsClientMock *mocks.MockSQSClient
	sqsAdapter    queue.SQSAdapter
}

func (s *sqsHandleTest) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.ctx = context.Background()
	s.sqsClientMock = mocks.NewMockSQSClient(ctrl)
	s.sqsAdapter = *queue.NewSQSAdapter(s.sqsClientMock, "errorQueueURL", "workQueueURL")
}

func Test_SQSHandleTest(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(sqsHandleTest))
}

func (s *sqsHandleTest) Test_Publish() {
	st := s.T()
	event := model.JobErrorEvent{
		JobID: "job-123",
		Email: "teste@gmail.com",
	}

	st.Run("should publish message successfully", func(t *testing.T) {
		expectedBody, err := json.Marshal(event)
		s.NoError(err)

		s.sqsClientMock.EXPECT().
			SendMessage(s.ctx, gomock.AssignableToTypeOf(&sqs.SendMessageInput{})).
			DoAndReturn(func(ctx context.Context, input *sqs.SendMessageInput, opts ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
				s.Equal("errorQueueURL", *input.QueueUrl)
				s.Equal(string(expectedBody), *input.MessageBody)
				return &sqs.SendMessageOutput{}, nil
			})

		err = s.sqsAdapter.Publish(s.ctx, event)
		s.NoError(err)
	})

	st.Run("should return error if SendMessage fails", func(t *testing.T) {

		s.sqsClientMock.EXPECT().
			SendMessage(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("send error"))

		err := s.sqsAdapter.Publish(s.ctx, event)
		s.Error(err)
		s.Contains(err.Error(), "failed to publish message to SQS")
	})
}

func (s *sqsHandleTest) Test_Receive() {
	st := s.T()

	st.Run("should return messages when SQS returns them", func(t *testing.T) {
		expectedMessages := []types.Message{
			{
				MessageId:     aws.String("msg-1"),
				ReceiptHandle: aws.String("handle-1"),
				Body:          aws.String("body-1"),
			},
			{
				MessageId:     aws.String("msg-2"),
				ReceiptHandle: aws.String("handle-2"),
				Body:          aws.String("body-2"),
			},
		}

		s.sqsClientMock.EXPECT().
			ReceiveMessage(gomock.Any(), gomock.Any()).
			Return(&sqs.ReceiveMessageOutput{Messages: expectedMessages}, nil)

		messages, err := s.sqsAdapter.Receive(s.ctx, 2, 1)

		s.NoError(err)
		s.Len(messages, 2)
		s.Equal(*expectedMessages[0].MessageId, *messages[0].MessageId)
		s.Equal(*expectedMessages[1].MessageId, *messages[1].MessageId)
	})

	st.Run("should return error when SQS returns error", func(t *testing.T) {
		expectedErr := fmt.Errorf("sqs error")
		s.sqsClientMock.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		messages, err := s.sqsAdapter.Receive(s.ctx, 5, 2)

		s.Error(err)
		s.Contains(err.Error(), "failed to receive messages from SQS")
		s.Nil(messages)
	})
}

func (s *sqsHandleTest) Test_Delete() {
	st := s.T()

	st.Run("should delete message successfully", func(t *testing.T) {
		receiptHandle := "test-receipt-handle"

		s.sqsClientMock.EXPECT().
			DeleteMessage(s.ctx, gomock.AssignableToTypeOf(&sqs.DeleteMessageInput{})).
			DoAndReturn(func(ctx context.Context, input *sqs.DeleteMessageInput, opts ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
				s.Equal("workQueueURL", *input.QueueUrl)
				s.Equal(receiptHandle, *input.ReceiptHandle)
				return &sqs.DeleteMessageOutput{}, nil
			})

		err := s.sqsAdapter.Delete(s.ctx, receiptHandle)
		s.NoError(err)
	})

	st.Run("should return error if DeleteMessage fails", func(t *testing.T) {
		receiptHandle := "fail-receipt-handle"

		s.sqsClientMock.EXPECT().
			DeleteMessage(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("delete error"))

		err := s.sqsAdapter.Delete(s.ctx, receiptHandle)
		s.Error(err)
		s.Contains(err.Error(), "failed to delete message from SQS")
	})
}

package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	model "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports"
)

type SQSAdapter struct {
	client        ports.SQSClient
	queueURL      string
	errorQueueURL string
}

func NewSQSAdapter(client ports.SQSClient, errorQueueURL, workQueueURL string) *SQSAdapter {
	return &SQSAdapter{
		client:        client,
		errorQueueURL: errorQueueURL,
		queueURL:      workQueueURL,
	}
}

func (s *SQSAdapter) Publish(ctx context.Context, event model.JobErrorEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize error message to JSON: %w", err)
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.errorQueueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to SQS: %w", err)
	}

	return nil
}

func (s *SQSAdapter) Receive(ctx context.Context, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	out, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTimeSeconds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from SQS: %w", err)
	}

	return out.Messages, nil
}

func (s *SQSAdapter) Delete(ctx context.Context, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from SQS: %w", err)
	}

	return nil
}

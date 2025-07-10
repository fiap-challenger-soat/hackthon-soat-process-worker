package input

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	model "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/domain"
	"github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/core/ports"
)

type Consumer struct {
	queue     ports.SQSAdapter
	processor ports.JobService
}

func NewConsumer(queue ports.SQSAdapter, processor ports.JobService) *Consumer {
	return &Consumer{
		queue:     queue,
		processor: processor,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("INFO: SQS consumer started. Listening for jobs...")
	for {
		select {
		case <-ctx.Done():
			log.Println("INFO: Consumer context cancelled. Exiting loop.")
			return
		default:
			msgs, err := c.queue.Receive(ctx, 10, 20)
			if err != nil {
				log.Printf("ERROR: Failed to receive messages: %v. Retrying in 10s...", err)
				time.Sleep(10 * time.Second)
				continue
			}
			for _, msg := range msgs {
				go c.handleMessage(context.Background(), msg)
			}
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg types.Message) {
	if msg.Body == nil || msg.ReceiptHandle == nil {
		log.Println("ERROR: Invalid message (nil body or receipt handle).")
		return
	}

	var jobMsg model.JobMessageEvent

	if err := json.Unmarshal([]byte(*msg.Body), &jobMsg); err != nil {
		log.Printf("ERROR: Failed to decode message body. Deleting message. Error: %v", err)
		_ = c.queue.Delete(ctx, *msg.ReceiptHandle)
		return
	}

	log.Printf("INFO: [Job %s] Processing started.", jobMsg.JobID)
	err := c.processor.ProcessJob(ctx, jobMsg.JobID, jobMsg.VideoPath)
	if err != nil {
		if delErr := c.queue.Delete(ctx, *msg.ReceiptHandle); delErr != nil {
			log.Printf("ERROR: [Job %s] Failed to delete message after error: %v", jobMsg.JobID, delErr)
		} else {
			log.Printf("INFO: [Job %s] ProcessJob failed. Message deleted", jobMsg.JobID)
		}
		return
	}

	if err := c.queue.Delete(ctx, *msg.ReceiptHandle); err != nil {
		log.Printf("ERROR: [Job %s] Failed to delete message: %v", jobMsg.JobID, err)
	} else {
		log.Printf("INFO: [Job %s] Message processed and deleted.", jobMsg.JobID)
	}
}

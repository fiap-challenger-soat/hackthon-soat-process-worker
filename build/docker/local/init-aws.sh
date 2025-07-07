#!/bin/bash
set -e

echo "=== [AWS Local Init] Initializing resources in LocalStack ==="

# Configuration
AWS_ENDPOINT_URL="http://localhost:4566"
BUCKET_INPUT="bucket-videos-entrada"
BUCKET_OUTPUT="bucket-videos-saida"
QUEUE_WORK="work-queue"
QUEUE_ERROR="error-queue"
VIDEO_FILE="/docker-entrypoint-initdb.d/name-of-your-file"  # <- Change this to your video file name
S3_KEY="uploads/video-inicial.mp4"
JOB_ID="194f2506-3a19-42fb-91a0-50442a1bfcfd"

# 1. Create S3 Buckets
echo "[1/6] Creating input S3 bucket: $BUCKET_INPUT"
aws --endpoint-url="$AWS_ENDPOINT_URL" s3 mb "s3://$BUCKET_INPUT"

echo "[2/6] Creating output S3 bucket: $BUCKET_OUTPUT"
aws --endpoint-url="$AWS_ENDPOINT_URL" s3 mb "s3://$BUCKET_OUTPUT"

# 2. Create SQS Queues
echo "[3/6] Creating work SQS queue: $QUEUE_WORK"
WORK_QUEUE_URL=$(aws --endpoint-url="$AWS_ENDPOINT_URL" sqs create-queue --queue-name "$QUEUE_WORK" --query 'QueueUrl' --output text)

echo "[4/6] Creating error SQS queue: $QUEUE_ERROR"
aws --endpoint-url="$AWS_ENDPOINT_URL" sqs create-queue --queue-name "$QUEUE_ERROR"

# 3. Upload Example Video
echo "[5/6] Uploading test video to S3: $VIDEO_FILE -> s3://$BUCKET_INPUT/$S3_KEY"
aws --endpoint-url="$AWS_ENDPOINT_URL" s3 cp "$VIDEO_FILE" "s3://$BUCKET_INPUT/$S3_KEY"

# 4. Send Initial Message to Work Queue
MESSAGE_BODY=$(printf '{"job_id": "%s", "video_path": "%s"}' "$JOB_ID" "$S3_KEY")
echo "[6/6] Sending initial message to SQS queue: $QUEUE_WORK"
aws --endpoint-url="$AWS_ENDPOINT_URL" sqs send-message \
  --queue-url "$WORK_QUEUE_URL" \
  --message-body "$MESSAGE_BODY"

echo "=== ✅ [AWS Local Init] Resources created and populated successfully! ==="

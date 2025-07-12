#!/bin/bash
set -e

echo "=== [AWS Local Init] Initializing resources in LocalStack ==="

# Configuration
AWS_ENDPOINT_URL="http://localhost:4566"
BUCKET_NAME="bucket-videos"
QUEUE_WORK="work-queue"
QUEUE_ERROR="error-queue"

# Arrays com os dados dos vídeos
VIDEO_FILES=(
  "/docker-entrypoint-initdb.d/your-video-file" # <- Add the name of your video file here
  "/docker-entrypoint-initdb.d/your-video-file" # <- Add the name of your video file here
)
S3_KEYS=(
  "uploads/video-inicial.mp4"
  "uploads/foto.jpg"
)
JOB_IDS=(
  "194f2506-3a19-42fb-91a0-50442a1bfcfd"
  "33daf232-990b-4411-8146-c5cd7c2e5c86"
)

# 1. Create S3 Bucket
echo "[1/5] Creating S3 bucket: $BUCKET_NAME"
aws --endpoint-url="http://localhost:4566" s3 mb "s3://bucket-videos"

# 2. Create SQS Queues
echo "[2/5] Creating work SQS queue: $QUEUE_WORK"
WORK_QUEUE_URL=$(aws --endpoint-url="http://localhost:4566" sqs create-queue --queue-name "work-queue" --query 'QueueUrl' --output text)

echo "[3/5] Creating error SQS queue: $QUEUE_ERROR"
aws --endpoint-url="http://localhost:4566" sqs create-queue --queue-name "error-queue"

# 3. Upload Example Videos
echo "[4/5] Uploading test videos to S3"
for i in "${!VIDEO_FILES[@]}"; do
  aws --endpoint-url="$AWS_ENDPOINT_URL" s3 cp "${VIDEO_FILES[$i]}" "s3://$BUCKET_NAME/${S3_KEYS[$i]}"
done

# 4. Send Initial Messages to Work Queue
echo "[5/5] Sending initial messages to SQS queue: $QUEUE_WORK"
for i in "${!VIDEO_FILES[@]}"; do
  MESSAGE_BODY=$(printf '{"job_id": "%s", "video_path": "%s"}' "${JOB_IDS[$i]}" "${S3_KEYS[$i]}")
  aws --endpoint-url="$AWS_ENDPOINT_URL" sqs send-message \
    --queue-url "$WORK_QUEUE_URL" \
    --message-body "$MESSAGE_BODY"
done

echo "=== ✅ [AWS Local Init] Resources created and populated successfully! ==="

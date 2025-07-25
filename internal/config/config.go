package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

var Vars appConfig

type appConfig struct {
	// DB config
	DBHost         string `env:"DB_HOST,required"`
	DBPort         int    `env:"DB_PORT,required"`
	DBUser         string `env:"DB_USER,required"`
	DBPassword     string `env:"DB_PASSWORD,required"`
	DBName         string `env:"DB_NAME,required"`
	DbMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS,required"`
	DbMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS,required"`

	// AWS config
	AWSRegion          string `env:"AWS_REGION,required"`
	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	AWSSessionToken    string `env:"AWS_SESSION,required"`
	AWSEndpointURL     string `env:"AWS_ENDPOINT_URL,required"`

	// S3 config
	S3Bucket string `env:"S3_BUCKET,required"`

	// SQS config
	SQSWorkQueueURL  string `env:"SQS_WORK_QUEUE_URL,required"`
	SQSErrorQueueURL string `env:"SQS_ERROR_QUEUE_URL,required"`
}

func Init() {
	if err := env.Parse(&Vars); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
}

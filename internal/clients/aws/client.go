package aws

import (
	"context"

	cgf "github.com/fiap-challenger-soat/hackthon-soat-process-worker/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewAwsConfig(ctx context.Context) (aws.Config, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cgf.Vars.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cgf.Vars.AWSAccessKeyID,
				cgf.Vars.AWSSecretAccessKey,
				cgf.Vars.AWSSessionToken,
			),
		),
	}

	if cgf.Vars.AWSEndpointURL != "" {
		opts = append(opts, config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: cgf.Vars.AWSEndpointURL,
				}, nil
			}),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func NewS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func NewSQSClient(cfg aws.Config) *sqs.Client {
	return sqs.NewFromConfig(cfg)
}

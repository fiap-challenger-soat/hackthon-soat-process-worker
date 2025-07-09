package aws

import (
	"context"

	cgf "github.com/fiap-challenger-soat/hackthon-soat-process-worker/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func NewAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cgf.Vars.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cgf.Vars.AWSAccessKeyID,
				cgf.Vars.AWSSecretAccessKey,
				cgf.Vars.AWSSessionToken,
			),
		),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: cgf.Vars.AWSEndpointURL}, nil
			}),
		),
	)
}

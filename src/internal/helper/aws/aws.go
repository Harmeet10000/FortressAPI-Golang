package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
)

type AWS struct {
	S3 *S3Client
}

func NewAWS(server *app.Server) (*AWS, error) {
	awsConfig := server.Config.S3

	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AccessKey,
			awsConfig.SecretKey,
			"",
		)),
	}

	// Add custom endpoint if provided (for S3-compatible services like Sevalla)
	if awsConfig.EndpointURL != "" {
		configOptions = append(configOptions, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string,
				options ...interface{},
			) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           awsConfig.EndpointURL,
					SigningRegion: awsConfig.Region,
				}, nil
			}),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, err
	}

	return &AWS{
		S3: NewS3Client(server, cfg),
	}, nil
}

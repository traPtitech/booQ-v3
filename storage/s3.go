package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/traPtitech/booQ-v3/domain"
)

type S3Storage struct {
	bucket string
	client *s3.Client
}

// NewS3Storage S3 (互換) オブジェクトストレージを作成します
func NewS3Storage(bucket, region, endpoint, apiKey, apiSecret string) (domain.FileStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(apiKey, apiSecret, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = &endpoint
	})

	return &S3Storage{
		bucket: bucket,
		client: client,
	}, nil
}

func (s *S3Storage) Save(filename string, src io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   src,
	}

	_, err := s.client.PutObject(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

func (s *S3Storage) Open(filename string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	output, err := s.client.GetObject(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return output.Body, nil
}

func (s *S3Storage) Delete(filename string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	_, err := s.client.DeleteObject(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

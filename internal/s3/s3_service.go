package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service provides methods to interact with an S3-compatible storage service.
type S3Service struct {
	Client *s3.Client
}

// NewService creates a new S3Service instance with the provided S3 client.
func NewService(cfg aws.Config) *S3Service {
	return &S3Service{
		Client: s3.NewFromConfig(cfg),
	}
}

// GetObject retrieves an object from the specified S3 bucket by its key.
func (s *S3Service) GetObject(bucket, key string) (io.ReadCloser, error) {
	req := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := s.Client.GetObject(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3 bucket %s with key %s: %w", bucket, key, err)
	}

	return resp.Body, nil
}

func (s *S3Service) ListObjectsForNotPrefix(bucket, notPrefix string) ([]string, error) {
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(notPrefix),
	}

	resp, err := s.Client.ListObjectsV2(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in S3 bucket %s with prefix %s: %w", bucket, notPrefix, err)
	}

	var keys []string
	for _, obj := range resp.Contents {
		if *obj.Key != notPrefix {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}

// ListObjectsForPrefix lists all objects in the specified S3 bucket that match the given prefix.
func (s *S3Service) ListObjectsForPrefix(bucket, prefix string) ([]string, error) {
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	resp, err := s.Client.ListObjectsV2(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in S3 bucket %s with prefix %s: %w", bucket, prefix, err)
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

// PutObject uploads a new object to the specified S3 bucket.
func (s *S3Service) PutObject(bucket, key, contentType string, body io.Reader) error {
	req := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        body,
	}

	_, err := s.Client.PutObject(context.Background(), req)
	return err
}

// DeleteObject removes an object from the specified S3 bucket.
func (s *S3Service) DeleteObject(bucket, key string) error {
	req := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.Client.DeleteObject(context.Background(), req)
	return err
}

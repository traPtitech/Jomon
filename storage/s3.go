package storage

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3 S3オブジェクトストレージ
type S3 struct {
	bucket string
	client *s3.Client
}

// SetS3Storage S3オブジェクトストレージをカレントストレージに設定します
func SetS3Storage(bucket string, region string, endpoint string, apiKey string, apiSecret string) (*S3, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service string, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				},
			),
		),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(apiKey, apiSecret, "")),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3{
		bucket: bucket,
		client: client,
	}, nil
}

func (s *S3) Save(filename string, src io.Reader) error {

	input := s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   src,
	}

	_, err := s.client.PutObject(context.Background(), &input)
	return err
}

func (s *S3) Open(filename string) (io.ReadCloser, error) {
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	res, err := s.client.GetObject(context.Background(), &input)
	if err != nil {
		var nf *types.NoSuchKey
		if errors.As(err, &nf) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return res.Body, nil
}

func (s *S3) Delete(filename string) error {
	input := s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}

	_, err := s.client.DeleteObject(context.Background(), &input)
	if err != nil {
		var nf *types.NoSuchKey
		if errors.As(err, &nf) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}

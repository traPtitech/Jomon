package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3FileStorage struct {
	bucket string
	client *s3.Client
}

func NewS3FileStorage(
	bucket string, region string, apiKey string, apiSecret string,
) (*S3FileStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(apiKey, apiSecret, ""),
		),
	)

	if err != nil {
		return &S3FileStorage{}, err
	}

	client := s3.NewFromConfig(cfg)

	s := &S3FileStorage{
		bucket: bucket,
		client: client,
	}

	return s, nil
}

func (fs *S3FileStorage) Save(filename string, key string, src io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket:             aws.String(fs.bucket),
		Key:                aws.String(key),
		Body:               src,
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}

	uploader := manager.NewUploader(fs.client)
	_, err := uploader.Upload(context.Background(), input)

	return err
}

func (fs *S3FileStorage) Open(key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}
	file, err := fs.client.GetObject(context.Background(), input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return file.Body, nil
}

func (fs *S3FileStorage) Delete(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}

	_, err := fs.client.DeleteObject(context.Background(), input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}

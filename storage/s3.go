package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	bucket string
	client *s3.Client
}

func NewS3Storage(cfg aws.Config, bucket string) (*S3Storage) {
	client := s3.NewFromConfig(cfg)

	s := &S3Storage{
		bucket: bucket,
		client: client,
	}

	return s
}

func (fs *S3Storage) Save(ctx context.Context, filename string, key string, src io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
		Body:   src,
		// FIXME : filenameに`"`が含まれているとバグりそう
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}

	uploader := manager.NewUploader(fs.client)
	_, err := uploader.Upload(ctx, input)

	return err
}

func (fs *S3Storage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}
	file, err := fs.client.GetObject(ctx, input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return file.Body, nil
}

func (fs *S3Storage) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}

	_, err := fs.client.DeleteObject(ctx, input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}

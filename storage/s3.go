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

var _ Storage = (*S3Storage)(nil)

func NewS3Storage(cfg aws.Config, bucket string) (*S3Storage) {
	client := s3.NewFromConfig(cfg)

	s := &S3Storage{
		bucket: bucket,
		client: client,
	}

	return s
}

// FIXME : 引数に`key string`を受け取り、`Key`フィールドに渡したいがinterfaceの定義上一時的に`filename`で代用してる　あとで分ける
func (fs *S3Storage) Save(filename string, src io.Reader) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(filename),
		Body:   src,
		// FIXME : filenameに`"`が含まれているとバグりそう
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}

	uploader := manager.NewUploader(fs.client)
	_, err := uploader.Upload(context.Background(), input)

	return err
}

func (fs *S3Storage) Open(key string) (io.ReadCloser, error) {
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

func (fs *S3Storage) Delete(key string) error {
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

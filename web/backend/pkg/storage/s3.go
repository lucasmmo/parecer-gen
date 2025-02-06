package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"parecer-gen/pkg/file"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	BucketClient *s3.Client
}

func NewS3Client(awsEndpoint, awsRegion string) S3Client {

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))

	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})
	return S3Client{BucketClient: s3Client}
}

func (cli S3Client) UploadFile(file *file.File) error {
	if _, err := cli.BucketClient.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("pdf-bucket"),
		Key:    aws.String(file.Filename),
		Body:   file.Reader,
	}); err != nil {
		return fmt.Errorf("error uploading parecer to S3: %w", err)
	}

	return nil
}

func (cli S3Client) DownloadFile(filename string) (io.Reader, error) {
	out, err := cli.BucketClient.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("pdf-bucket"),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

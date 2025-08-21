package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MC *MinioClient
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

// مقداردهی مستقیم به MC
func InitMinio(endpoint, accessKey, secretKey, bucket string, useSSL bool) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to create MinIO client: %v", err)
	}

	ctx := context.Background()
	exists, errBucketExists := client.BucketExists(ctx, bucket)
	if errBucketExists != nil {
		log.Fatalf("error checking bucket: %v", errBucketExists)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
		log.Printf("✅ bucket %s created", bucket)
	} else {
		log.Printf("✅ bucket %s already exists", bucket)
	}

	MC = &MinioClient{
		Client: client,
		Bucket: bucket,
	}
}

func (mc *MinioClient) UploadFile(ctx context.Context, reader io.Reader, fileName, contentType string, fileSize int64) (string, error) {
	_, err := mc.Client.PutObject(ctx, mc.Bucket, fileName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url, err := mc.Client.PresignedGetObject(ctx, mc.Bucket, fileName, time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return url.String(), nil
}

func (mc *MinioClient) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	obj, err := mc.Client.GetObject(ctx, mc.Bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return obj, nil
}

func (mc *MinioClient) DeleteFile(ctx context.Context, fileName string) error {
	err := mc.Client.RemoveObject(ctx, mc.Bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

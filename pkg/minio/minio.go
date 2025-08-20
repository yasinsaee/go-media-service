package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) *MinioClient {
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

	return &MinioClient{
		Client: client,
		Bucket: bucket,
	}
}

// 🔹 Upload فایل
func (mc *MinioClient) UploadFile(ctx context.Context, file multipart.File, fileName, contentType string, fileSize int64) (string, error) {
	_, err := mc.Client.PutObject(ctx, mc.Bucket, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// لینک public presigned
	url, err := mc.Client.PresignedGetObject(ctx, mc.Bucket, fileName, time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return url.String(), nil
}

// 🔹 دانلود فایل
func (mc *MinioClient) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	obj, err := mc.Client.GetObject(ctx, mc.Bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return obj, nil
}

// 🔹 حذف فایل
func (mc *MinioClient) DeleteFile(ctx context.Context, fileName string) error {
	err := mc.Client.RemoveObject(ctx, mc.Bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

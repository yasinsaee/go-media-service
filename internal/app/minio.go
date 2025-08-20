package app

import "github.com/yasinsaee/go-media-service/pkg/minio"

func InitMinio() {
	minio.InitMinio(
		"localhost:9000",
		"admin",
		"password123",
		"media",
		false,
	)
}

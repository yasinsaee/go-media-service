package app

import (
	"strconv"

	"github.com/yasinsaee/go-media-service/internal/app/config"
	"github.com/yasinsaee/go-media-service/pkg/minio"
)

func initMinio() {
	ssl, _ := strconv.ParseBool(config.GetEnv("MINIO_SSL", "false"))

	minio.InitMinio(
		config.GetEnv("MINIO_URI", "localhost:9000"),
		config.GetEnv("MINIO_USERNAME", "admin"),
		config.GetEnv("MINIO_PASSWORD", "password123"),
		config.GetEnv("MINIO_BUCKET", "media"),
		ssl,
	)
}

package app

import "github.com/yasinsaee/go-media-service/internal/app/config"

func StartApp() {
	config.LoadEnv()

	InitMongo()

	InitMinio()

	StartGRPCServer()

}

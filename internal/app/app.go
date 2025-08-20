package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func StartApp() {
	config.LoadEnv()

	InitMongo()

	// InitJWT()

	StartGRPCServer()
	
}

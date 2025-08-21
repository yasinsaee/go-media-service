package app

import (
	"log"
	"net"

	mediagrpc "github.com/yasinsaee/go-media-service/internal/handlers/grpc/media"
	repository_media "github.com/yasinsaee/go-media-service/internal/repository/media"
	"github.com/yasinsaee/go-media-service/internal/service/media"
	"github.com/yasinsaee/go-media-service/pkg/mongo"
	mediapb "github.com/yasinsaee/go-media-service/user-service/media"

	"google.golang.org/grpc"
)

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// s := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor()))
	s := grpc.NewServer()

	//repos
	mediaRepo := repository_media.NewMongoMediaRepository(mongo.DB.Database, "media")

	//services
	mediaService := media.NewMediaService(mediaRepo)

	//handlers
	mediaHandler := mediagrpc.New(mediaService)

	//register grpc services
	mediapb.RegisterPermissionServiceServer(s, mediaHandler)

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

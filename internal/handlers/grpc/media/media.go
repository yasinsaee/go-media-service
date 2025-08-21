package media

import (
	"bytes"
	"context"

	mediapb "github.com/yasinsaee/go-media-service/media-service/media"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yasinsaee/go-media-service/internal/domain/media"
)

type MediaServiceServer struct {
	mediapb.UnimplementedMediaServiceServer
	service media.MediaService
}

// UploadMedia handler
// UploadMedia gRPC handler
func (s *MediaServiceServer) UploadMedia(ctx context.Context, req *mediapb.UploadMediaRequest) (*mediapb.UploadMediaResponse, error) {
	if req == nil || len(req.FileContent) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	fileReader := bytes.NewReader(req.FileContent)

	mediaObj := &media.Media{
		OriginalName: req.OriginalName,
		OwnerID:      req.OwnerId,
		Tags:         req.Tags,
		Metadata:     req.Metadata,
		Privacy:      req.Privacy,
		Status:       "active",
	}

	createdMedia, err := s.service.UploadMedia(ctx, mediaObj, fileReader, int64(len(req.FileContent)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "upload failed: %v", err)
	}

	return &mediapb.UploadMediaResponse{
		Media: &mediapb.Media{
			Id:           createdMedia.ID,
			FileName:     createdMedia.FileName,
			OriginalName: createdMedia.OriginalName,
			ContentType:  createdMedia.ContentType,
			Size:         createdMedia.Size,
			Url:          createdMedia.URL,
			Thumbnail:    createdMedia.Thumbnail,
			OwnerId:      createdMedia.OwnerID,
			Tags:         createdMedia.Tags,
			Metadata:     createdMedia.Metadata,
			Privacy:      createdMedia.Privacy,
			Status:       createdMedia.Status,
			ExpiresAt:    timestamppb.New(createdMedia.ExpiresAt),
			CreatedAt:    timestamppb.New(createdMedia.CreatedAt),
			UpdatedAt:    timestamppb.New(createdMedia.UpdatedAt),
		},
	}, nil
}

// func (s *MediaServiceServer) GetMedia(ctx context.Context, req *mediapb.GetMediaRequest) (*mediapb.MediaResponse, error) {
// 	mediaObj, err := s.Service.GetByID(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &mediapb.MediaResponse{
// 		Media: &mediapb.Media{
// 			Id:           mediaObj.ID,
// 			FileName:     mediaObj.FileName,
// 			OriginalName: mediaObj.OriginalName,
// 			ContentType:  mediaObj.ContentType,
// 			Size:         mediaObj.Size,
// 			Url:          mediaObj.URL,
// 			Thumbnail:    mediaObj.Thumbnail,
// 			OwnerId:      mediaObj.OwnerID,
// 			Tags:         mediaObj.Tags,
// 			Privacy:      mediaObj.Privacy,
// 			Status:       mediaObj.Status,
// 			ExpiresAt:    mediaObj.ExpiresAt.Format(time.RFC3339),
// 			CreatedAt:    mediaObj.CreatedAt.Format(time.RFC3339),
// 			UpdatedAt:    mediaObj.UpdatedAt.Format(time.RFC3339),
// 		},
// 	}, nil
// }

// func (s *MediaServiceServer) DeleteMedia(ctx context.Context, req *mediapb.DeleteMediaRequest) (*mediapb.DeleteMediaResponse, error) {
// 	err := s.Service.Delete(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &mediapb.DeleteMediaResponse{Success: true}, nil
// }

// func (s *MediaServiceServer) ListMedia(ctx context.Context, req *mediapb.ListMediaRequest) (*mediapb.ListMediaResponse, error) {
// 	medias, err := s.Service.List(ctx, req.OwnerId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var protoMedias []*mediapb.Media
// 	for _, m := range medias {
// 		protoMedias = append(protoMedias, &mediapb.Media{
// 			Id:           m.ID,
// 			FileName:     m.FileName,
// 			OriginalName: m.OriginalName,
// 			ContentType:  m.ContentType,
// 			Size:         m.Size,
// 			Url:          m.URL,
// 			Thumbnail:    m.Thumbnail,
// 			OwnerId:      m.OwnerID,
// 			Tags:         m.Tags,
// 			Privacy:      m.Privacy,
// 			Status:       m.Status,
// 			ExpiresAt:    m.ExpiresAt.Format(time.RFC3339),
// 			CreatedAt:    m.CreatedAt.Format(time.RFC3339),
// 			UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
// 		})
// 	}

// 	return &mediapb.ListMediaResponse{Medias: protoMedias}, nil
// }

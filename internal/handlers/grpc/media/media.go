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

type Handler struct {
	mediapb.UnimplementedMediaServiceServer
	service media.MediaService
}

func New(service media.MediaService) *Handler {
	return &Handler{service: service}
}

// UploadMedia handler
// UploadMedia gRPC handler
func (s *Handler) UploadMedia(ctx context.Context, req *mediapb.UploadMediaRequest) (*mediapb.UploadMediaResponse, error) {
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

func (s *Handler) GetMedia(ctx context.Context, req *mediapb.GetMediaByIDRequest) (*mediapb.GetMediaByIDResponse, error) {
	mediaObj, err := s.service.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &mediapb.GetMediaByIDResponse{
		Media: &mediapb.Media{
			Id:           mediaObj.ID,
			FileName:     mediaObj.FileName,
			OriginalName: mediaObj.OriginalName,
			ContentType:  mediaObj.ContentType,
			Size:         mediaObj.Size,
			Url:          mediaObj.URL,
			Thumbnail:    mediaObj.Thumbnail,
			OwnerId:      mediaObj.OwnerID,
			Tags:         mediaObj.Tags,
			Privacy:      mediaObj.Privacy,
			Status:       mediaObj.Status,
			ExpiresAt:    timestamppb.New(mediaObj.ExpiresAt),
			CreatedAt:    timestamppb.New(mediaObj.CreatedAt),
			UpdatedAt:    timestamppb.New(mediaObj.UpdatedAt),
		},
	}, nil
}

func (s *Handler) DeleteMedia(ctx context.Context, req *mediapb.DeleteMediaRequest) (*mediapb.DeleteMediaResponse, error) {
	err := s.service.Delete(req.Id, req.Force)
	if err != nil {
		return nil, err
	}
	return &mediapb.DeleteMediaResponse{Success: true}, nil
}

func (s *Handler) ListMedia(ctx context.Context, req *mediapb.ListMediaRequest) (*mediapb.ListMediaResponse, error) {
	filter := make(map[string]any, len(req.Filter))
	for k, v := range req.Filter {
		filter[k] = v
	}

	medias, err := s.service.List(filter, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var protoMedias []*mediapb.Media
	for _, m := range medias {
		protoMedias = append(protoMedias, &mediapb.Media{
			Id:           m.ID,
			FileName:     m.FileName,
			OriginalName: m.OriginalName,
			ContentType:  m.ContentType,
			Size:         m.Size,
			Url:          m.URL,
			Thumbnail:    m.Thumbnail,
			OwnerId:      m.OwnerID,
			Tags:         m.Tags,
			Privacy:      m.Privacy,
			Status:       m.Status,
			ExpiresAt:    timestamppb.New(m.ExpiresAt),
			CreatedAt:    timestamppb.New(m.CreatedAt),
			UpdatedAt:    timestamppb.New(m.UpdatedAt),
		})
	}

	return &mediapb.ListMediaResponse{Medias: protoMedias}, nil
}


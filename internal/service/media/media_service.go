package media

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/yasinsaee/go-media-service/internal/domain/media"
	"github.com/yasinsaee/go-media-service/pkg/minio"
)

type mediaService struct {
	repo media.MediaRepository
}

// NewMediaService creates a new MediaService
func NewMediaService(repo media.MediaRepository) media.MediaService {
	return &mediaService{
		repo: repo,
	}
}

func (s *mediaService) Create(m *media.Media) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
	return s.repo.Create(m)
}

func (s *mediaService) GetByID(id any) (*media.Media, error) {
	return s.repo.GetByID(id)
}

func (s *mediaService) Update(m *media.Media) error {
	return s.repo.Update(m)
}

func (s *mediaService) Delete(id any, force bool) error {
	mediaItem, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if mediaItem != nil {
		_ = minio.MC.DeleteFile(context.Background(), mediaItem.FileName)
	}

	return s.repo.Delete(id, force)
}

func (s *mediaService) List(filter map[string]any, limit, offset int) (media.Medias, error) {
	return s.repo.List(filter, limit, offset)
}

func (s *mediaService) GetByOwner(ownerID string, limit, offset int) (media.Medias, error) {
	return s.repo.GetByOwner(ownerID, limit, offset)
}

func (s *mediaService) GetByTags(tags []string, limit, offset int) (media.Medias, error) {
	return s.repo.GetByTags(tags, limit, offset)
}

func (s *mediaService) GetExpired() (media.Medias, error) {
	return s.repo.GetExpired()
}

func (s *mediaService) UploadMedia(ctx context.Context, m *media.Media, file multipart.File, fileSize int64) (*media.Media, error) {
	url, err := minio.MC.UploadFile(ctx, file, m.FileName, m.ContentType, fileSize)
	if err != nil {
		return nil, err
	}

	m.URL = url
	if err := s.repo.Create(m); err != nil {
		_ = minio.MC.DeleteFile(ctx, m.FileName)
		return nil, err
	}

	return m, nil
}

func (s *mediaService) GetFile(ctx context.Context, fileName string) (interface{}, error) {
	return minio.MC.GetFile(ctx, fileName)
}

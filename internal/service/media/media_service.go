package media

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
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

func (s *mediaService) UploadMedia(ctx context.Context, m *media.Media, file io.Reader, fileSize int64) (*media.Media, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	m.ContentType = http.DetectContentType(fileBytes[:512])

	ext := filepath.Ext(m.OriginalName)
	if ext == "" {
		if exts, _ := mime.ExtensionsByType(m.ContentType); len(exts) > 0 {
			ext = exts[0]
		}
	}
	m.FileName = uuid.New().String() + ext
	m.Size = int64(len(fileBytes))

	url, err := minio.MC.UploadFile(ctx, bytes.NewReader(fileBytes), m.FileName, m.ContentType, m.Size)
	if err != nil {
		return nil, err
	}
	m.URL = url

	if strings.HasPrefix(m.ContentType, "image/") {
		thumbURL, err := s.GenerateThumbnail(ctx, bytes.NewReader(fileBytes), m.FileName)
		if err == nil {
			m.Thumbnail = thumbURL
		}
	}

	if err := s.repo.Create(m); err != nil {
		_ = minio.MC.DeleteFile(ctx, m.FileName)
		return nil, err
	}

	return m, nil
}

func (s *mediaService) GetFile(ctx context.Context, fileName string) (interface{}, error) {
	return minio.MC.GetFile(ctx, fileName)
}

func (s *mediaService) GenerateThumbnail(ctx context.Context, file io.Reader, fileName string) (string, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, thumb, nil); err != nil {
		return "", err
	}

	thumbName := "thumb_" + fileName
	url, err := minio.MC.UploadFile(ctx, buf, thumbName, "image/jpeg", int64(buf.Len()))
	if err != nil {
		return "", err
	}

	return url, nil
}

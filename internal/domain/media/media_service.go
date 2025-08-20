package media

import (
	"context"
	"mime/multipart"
)

type MediaService interface {
	Create(media *Media) error
	GetByID(id any) (*Media, error)
	Update(media *Media) error
	Delete(id any, force bool) error
	List(filter map[string]any, limit, offset int) (Medias, error)
	GetByOwner(ownerID string, limit, offset int) (Medias, error)
	GetByTags(tags []string, limit, offset int) (Medias, error)
	GetExpired() (Medias, error)
	UploadMedia(ctx context.Context, m *Media, file multipart.File, fileSize int64) (*Media, error)
	GetFile(ctx context.Context, fileName string) (interface{}, error)
}

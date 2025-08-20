package media

import "github.com/yasinsaee/go-media-service/internal/domain/media"

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
	return s.repo.Create(m)
}

func (s *mediaService) GetByID(id any) (*media.Media, error) {
	return s.repo.GetByID(id)
}

func (s *mediaService) Update(m *media.Media) error {
	return s.repo.Update(m)
}

func (s *mediaService) Delete(id any, force bool) error {
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

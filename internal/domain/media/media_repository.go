package media

type MediaRepository interface {
	Create(m *Media) error
	GetByID(id any) (*Media, error)
	Update(m *Media) error
	Delete(id any, force bool) error
	List(filter map[string]any, limit, offset int) (Medias, error)
	GetByOwner(ownerID string, limit, offset int) (Medias, error)
	GetByTags(tags []string, limit, offset int) (Medias, error)
	GetExpired() (Medias, error)
}

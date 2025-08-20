package repository

import (
	"time"

	"github.com/yasinsaee/go-media-service/internal/domain/media"
	"github.com/yasinsaee/go-media-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-media-service/pkg/mongo"
	"github.com/yasinsaee/go-media-service/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoMediaRepository implements the MediaRepository interface using MongoDB.
type mongoMediaRepository struct {
	collection *mongo.Collection
}

// NewMongoMediaRepository returns a new instance of mongoMediaRepository.
func NewMongoMediaRepository(db *mongo.Database, collectionName string) media.MediaRepository {
	return &mongoMediaRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new media document into the database.
func (r *mongoMediaRepository) Create(m *media.Media) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return mongo2.Create(m)
}

// GetByID retrieves a media document by ID.
func (r *mongoMediaRepository) GetByID(id any) (*media.Media, error) {
	m := new(media.Media)
	err := mongo2.Get(r.collection.Name(), id, m)
	return m, err
}

// Update modifies an existing media document.
func (r *mongoMediaRepository) Update(m *media.Media) error {
	m.UpdatedAt = time.Now()
	return mongo2.Update(m)
}

// Delete removes a media document by ID. If force is false, it can do soft delete.
func (r *mongoMediaRepository) Delete(id any, force bool) error {
	if !force {
		// soft delete
		m, err := r.GetByID(id)
		if err != nil {
			return err
		}
		m.DeletedAt = time.Now()
		return mongo2.Update(m)
	}

	// hard delete
	objID, err := util.ToObjectID(id)
	if err != nil {
		logger.Error("error while deleting media: ", err.Error())
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

// List returns media documents with optional filters and pagination.
func (r *mongoMediaRepository) List(filter map[string]any, limit, offset int) (media.Medias, error) {
	medias := make(media.Medias, 0)
	err := mongo2.Find(r.collection.Name(), filter, limit, offset, &medias)
	if err != nil {
		logger.Error("error while fetching media: ", err.Error())
		return nil, err
	}
	return medias, nil
}

// GetByOwner returns media documents owned by a specific owner.
func (r *mongoMediaRepository) GetByOwner(ownerID string, limit, offset int) (media.Medias, error) {
	filter := bson.M{"owner_id": ownerID}
	return r.List(filter, limit, offset)
}

// GetByTags returns media documents matching given tags.
func (r *mongoMediaRepository) GetByTags(tags []string, limit, offset int) (media.Medias, error) {
	filter := bson.M{"tags": bson.M{"$in": tags}}
	return r.List(filter, limit, offset)
}

// GetExpired returns media documents that are expired.
func (r *mongoMediaRepository) GetExpired() (media.Medias, error) {
	medias := make(media.Medias, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{"expires_at": bson.M{"$lte": time.Now()}}, &medias)
	if err != nil {
		logger.Error("error while fetching expired media: ", err.Error())
		return nil, err
	}
	return medias, nil
}

package media

import "time"

type Media struct {
	ID           string            `json:"id" bson:"_id,omitempty"`
	FileName     string            `json:"file_name" bson:"file_name"`
	OriginalName string            `json:"original_name" bson:"original_name"`
	ContentType  string            `json:"content_type" bson:"content_type"`
	Size         int64             `json:"size" bson:"size"`
	URL          string            `json:"url" bson:"url"`
	Thumbnail    string            `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	OwnerID      string            `json:"owner_id" bson:"owner_id"`
	Tags         []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Privacy      string            `json:"privacy" bson:"privacy"` // public, private, internal
	Status       string            `json:"status" bson:"status"`   // active, deleted, archived
	ExpiresAt    time.Time         `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	CreatedAt    time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt    time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type Medias []Media

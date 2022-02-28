package mongol

import (
	"time"

	timecop "github.com/bluele/go-timecop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseDocument
type BaseDocument struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetID() returns ID of the document
func (m *BaseDocument) GetID() primitive.ObjectID {
	return m.ID
}

// GetHexID() returns ID of the document as a hex-string
func (m *BaseDocument) GetHexID() string {
	return m.ID.Hex()
}

// SetHexID() sets ID of the document from the hex-string
func (m *BaseDocument) SetHexID(hexID string) error {
	oid, err := StringToObjectID(hexID)
	if err != nil {
		return err
	}

	m.ID = oid

	return nil
}

// SetupCreatedAt() sets CreatedAt field for the model
// The method does not store any data to database
// you should use the method before InsertMany(), UpdateMany() requests from you storage
func (m *BaseDocument) SetupCreatedAt() {
	m.CreatedAt = timecop.Now().UTC()
}

// SetupUpdatedAt() sets UpdatedAt field for the model
// The method does not store any data to database
// you should use the method before InsertMany(), UpdateMany() requests from you storage
func (m *BaseDocument) SetupUpdatedAt() {
	m.UpdatedAt = timecop.Now().UTC()
}

// SetJSONID() sets the model ID from given JSON
func (m *BaseDocument) SetJSONID(jsonB []byte) error {
	return m.ID.UnmarshalJSON(jsonB)
}

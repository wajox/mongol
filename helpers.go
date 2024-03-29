package mongol

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDsFromStrArr() converts provided array of mongoid strings to ObjectID
func ObjectIDsFromStrArr(idsRaw []string) ([]*primitive.ObjectID, error) {
	ids := make([]*primitive.ObjectID, len(idsRaw))

	for i := 0; i < cap(ids); i++ {
		id, err := primitive.ObjectIDFromHex(idsRaw[i])
		if err != nil {
			return nil, fmt.Errorf("objectID conversion err: %s", err)
		}

		ids[i] = &id
	}

	return ids, nil
}

// IsValidMongoID() validates provided mongoid string
func IsValidMongoID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)

	return err == nil
}

// StringToObjectID() converts the string to Mongo ObjectID
func StringToObjectID(hexID string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return primitive.ObjectID{}, ErrInvalidObjectID
	}

	return oid, nil
}

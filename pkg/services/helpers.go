package services

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDsFromStrArr - converts provided array of mongoid strings to ObjectID
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

// IsValidMongoID - validate provided mongoid string
func IsValidMongoID(id string) bool {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return false
	}

	return true
}

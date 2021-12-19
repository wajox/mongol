package helpers

import (
	bs "github.com/wajox/mongol/pkg/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObjectID(hexID string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return primitive.ObjectID{}, bs.ErrInvalidObjectID
	}

	return oid, nil
}

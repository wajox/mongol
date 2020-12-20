package services

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// numeric code for MongoDB duplication error
	DuplicationErrorCode = 11000
)

var (
	// ErrDocumentNotFound appears then the document was not found in the collection
	ErrDocumentNotFound = errors.New("can not find document")
	// ErrDocumentDuplication appears then the document has fields that already exists in another document
	ErrDocumentDuplication = errors.New("documents duplication was detected")
	// ErrDocumentNotModified appears then the document was not modified after your request
	ErrDocumentNotModified = errors.New("document wasn't modified")
	// ErrInvalidObjectID appears then the ID has invalid format
	ErrInvalidObjectID = errors.New("invalid objectID")
)

// HandleDuplicationErr() checks exception type
// if the error has occurred due to duplication problem then
// the method returns ErrDocumentDuplication
func HandleDuplicationErr(err error) error {
	mngErr, isWriteErr := err.(mongo.WriteException)
	if !isWriteErr || mngErr.WriteErrors[0].Code != DuplicationErrorCode {
		return err
	}

	return ErrDocumentDuplication
}

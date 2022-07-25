package mongol_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/wajox/mongol"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ = Describe("Errors", func() {
	Describe("HandleDuplicationErr()", func() {
		It("should return duplication error", func() {
			mongoErr := mongo.WriteException{
				WriteConcernError: nil,
				WriteErrors: mongo.WriteErrors{
					{
						Code: mongol.DuplicationErrorCode,
					},
				},
			}

			Expect(mongol.HandleDuplicationErr(mongoErr)).To(Equal(mongol.ErrDocumentDuplication))
		})

		It("should return the original error", func() {
			mongoErr := errors.New("some error")

			Expect(mongol.HandleDuplicationErr(mongoErr)).To(Equal(mongoErr))
		})
	})
})

package mongol_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wajox/mongol"
)

var _ = Describe("Helpers", func() {
	ValidObjectID := "6291dc08a802d7000622f16a"
	InvalidObjectID := "6291dc08a802d7000622f"

	Describe("IsValidMongoId()", func() {
		It("should validate mongo ObjectID", func() {
			Expect(mongol.IsValidMongoID(ValidObjectID)).To(BeTrue())
			Expect(mongol.IsValidMongoID(InvalidObjectID)).To(BeFalse())
		})
	})

	Describe("ObjectIDsFromStrArr()", func() {
		It("should convert strings to ObjectIDs", func() {
			valid, errValid := mongol.ObjectIDsFromStrArr([]string{ValidObjectID})

			Expect(valid[0].Hex()).To(Equal(ValidObjectID))
			Expect(errValid).To(BeNil())

			invalid, errInvalid := mongol.ObjectIDsFromStrArr([]string{InvalidObjectID})

			Expect(invalid).To(BeNil())
			Expect(errInvalid).NotTo(BeNil())
		})
	})

	Describe("StringToObjectID()", func() {
		It("should convert string to object ID", func() {
			validOid, validErr := mongol.StringToObjectID(ValidObjectID)
			Expect(validOid.Hex()).To(Equal(ValidObjectID))
			Expect(validErr).To(BeNil())

			invalidOid, invalidErr := mongol.StringToObjectID(InvalidObjectID)
			Expect(invalidOid.Hex()).To(Equal("000000000000000000000000"))
			Expect(invalidErr).NotTo(BeNil())
		})
	})
})

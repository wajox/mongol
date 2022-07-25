package mongol_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/wajox/mongol"
)

var _ = Describe("BaseDocument", func() {
	Describe("SetHexID()", func() {
		ValidObjectID := "6291dc08a802d7000622f16a"
		InvalidObjectID := "6291dc08a802d7000622f"

		It("should set an ID from string", func() {
			doc := &mongol.BaseDocument{}

			Expect(doc.SetHexID(ValidObjectID)).To(BeNil())
			Expect(doc.GetHexID()).To(Equal(ValidObjectID))
		})

		It("should not set an ID from string", func() {
			doc := &mongol.BaseDocument{}

			Expect(doc.SetHexID(InvalidObjectID)).NotTo(BeNil())
			Expect(doc.GetHexID()).To(Equal("000000000000000000000000"))
		})
	})
})

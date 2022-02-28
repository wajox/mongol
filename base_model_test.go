package mongol_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	timecop "github.com/bluele/go-timecop"
	. "github.com/wajox/mongol"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("BaseDocument", func() {
	Context("with existing model", func() {
		var (
			m        *BaseDocument
			objectID primitive.ObjectID
		)
		BeforeEach(func() {
			objectID = primitive.NewObjectID()
			m = &BaseDocument{ID: objectID}
		})

		Describe("constructor", func() {
			It("should create new model", func() {
				Expect(m).NotTo(BeNil())
				Expect(m.ID).To(Equal(objectID))
			})
		})

		Describe("methods", func() {
			Describe("GetID()", func() {
				It("should return ObjectID", func() {
					Expect(m.GetID()).To(Equal(objectID))
				})
			})

			Describe("GetHexID()", func() {
				It("should return Hex ID", func() {
					Expect(m.GetHexID()).To(Equal(objectID.Hex()))
				})
			})

			Describe("SetHexID()", func() {
				It("should set Hex ID", func() {
					newObjectID := primitive.NewObjectID()
					err := m.SetHexID(newObjectID.Hex())

					Expect(err).To(BeNil())
					Expect(m.GetHexID()).NotTo(Equal(objectID.Hex()))
					Expect(m.GetHexID()).To(Equal(newObjectID.Hex()))
				})
			})

			Describe("SetupTimestamps()", func() {
				Context("when created_at and updated_at are empty", func() {
					It("should set created_at, updated_at to current time", func() {
						timecop.Freeze(time.Now().UTC().Add(time.Hour * 1))
						defer timecop.Return()

						m.SetupCreatedAt()
						m.SetupUpdatedAt()

						Expect(m.CreatedAt).To(Equal(timecop.Now().UTC()))
						Expect(m.UpdatedAt).To(Equal(timecop.Now().UTC()))
					})
				})

				Context("when created_at is not empty", func() {
					var (
						prevCreatedAt time.Time
					)
					BeforeEach(func() {
						prevCreatedAt = timecop.Now().UTC()
						m.CreatedAt = prevCreatedAt
						m.UpdatedAt = timecop.Now().UTC()
					})

					It("should update only updated_at field", func() {
						timecop.Freeze(time.Now().UTC().Add(time.Hour * 1))
						defer timecop.Return()

						m.SetupUpdatedAt()

						Expect(m.CreatedAt).To(Equal(prevCreatedAt))
						Expect(m.UpdatedAt).To(Equal(timecop.Now().UTC()))
					})
				})
			})
		})
	})
})

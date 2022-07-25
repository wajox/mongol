package mongol_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/wajox/mongol"
)

var _ = Describe("FilterBuilder", func() {
	Describe("NewFilterBuilder()", func() {
		It("should create a new builder", func() {
			Expect(mongol.NewFilterBuilder()).NotTo(BeNil())
		})
	})

	Describe("methods", func() {
		var (
			filter *mongol.FilterBuilder
		)

		Context("with new filter builder", func() {
			BeforeEach(func() {
				filter = mongol.NewFilterBuilder()
			})

			Describe("Where()", func() {
				It("should add a new condition", func() {
					filter.Where("name", bson.M{"$eq": "John"})

					query := filter.GetQuery()

					Expect(query["name"]).To(Equal(bson.M{"$eq": "John"}))
				})
			})

			Describe("Or()", func() {
				It("should add a new condition", func() {
					filter.Or(
						bson.M{"name": "John"},
						bson.M{"name": "Mike"},
					)

					query := filter.GetQuery()

					Expect(len(query["$or"].(bson.A))).To(Equal(2))
				})
			})

			Describe("Or()", func() {
				It("should add a new condition", func() {
					filter.And(
						bson.M{"name": "John"},
						bson.M{"last_name": "Doe"},
					)

					query := filter.GetQuery()

					Expect(len(query["$and"].(bson.A))).To(Equal(2))
				})
			})

			Describe("EqualTo()", func() {
				It("should add a new condition", func() {
					filter.EqualTo("name", "John")

					query := filter.GetQuery()

					Expect(query["name"]).To(Equal(bson.M{"$eq": "John"}))
				})
			})

			Describe("NotEqualTo()", func() {
				It("should add a new condition", func() {
					filter.NotEqualTo("name", "John")

					query := filter.GetQuery()

					Expect(query["name"]).To(Equal(bson.M{"$ne": "John"}))
				})
			})
		})
	})
})

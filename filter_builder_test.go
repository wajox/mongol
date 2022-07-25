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

			Describe("In()", func() {
				It("should add a new condition", func() {
					filter.In("name", bson.A{"John"})

					query := filter.GetQuery()

					Expect(query["name"].(bson.M)["$in"].(bson.A)).To(Equal(bson.A{"John"}))
				})
			})

			Describe("NotIn()", func() {
				It("should add a new condition", func() {
					filter.NotIn("name", bson.A{"John"})

					query := filter.GetQuery()

					Expect(query["name"].(bson.M)["$nin"].(bson.A)).To(Equal(bson.A{"John"}))
				})
			})

			Describe("HasField()", func() {
				It("should add a new condition", func() {
					filter.HasField("name")

					query := filter.GetQuery()

					Expect(query["name"].(bson.M)["$exists"].(bool)).To(BeTrue())
				})
			})

			Describe("HasNotField()", func() {
				It("should add a new condition", func() {
					filter.HasNotField("name")

					query := filter.GetQuery()

					Expect(query["name"].(bson.M)["$exists"].(bool)).To(BeFalse())
				})
			})

			Describe("Gte()", func() {
				It("should add a new condition", func() {
					filter.Gte("age", 10)

					query := filter.GetQuery()

					Expect(query["age"]).To(Equal(bson.M{"$gte": 10}))
				})
			})

			Describe("Lte()", func() {
				It("should add a new condition", func() {
					filter.Lte("age", 10)

					query := filter.GetQuery()

					Expect(query["age"]).To(Equal(bson.M{"$lte": 10}))
				})
			})

			Describe("Gt()", func() {
				It("should add a new condition", func() {
					filter.Gt("age", 10)

					query := filter.GetQuery()

					Expect(query["age"]).To(Equal(bson.M{"$gt": 10}))
				})
			})

			Describe("Lt()", func() {
				It("should add a new condition", func() {
					filter.Lt("age", 10)

					query := filter.GetQuery()

					Expect(query["age"]).To(Equal(bson.M{"$lt": 10}))
				})
			})
		})
	})
})

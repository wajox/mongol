package mongol_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	timecop "github.com/bluele/go-timecop"
	"github.com/google/uuid"
	. "github.com/wajox/mongol"
)

type ExampleModel struct {
	BaseDocument `bson:",inline"`

	Title string `json:"title,omitempty" bson:"title,omitempty"`
}

func NewExampleModel() *ExampleModel {
	return &ExampleModel{Title: "Test title " + uuid.New().String()}
}

// nolint
var _ = Describe("BaseCollection", func() {
	var (
		mongoURI, mongoDBName, mongoCollectionName string
		storage                                    *BaseCollection
		connErr                                    error
	)

	BeforeEach(func() {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://0.0.0.0:27017"
		}
		mongoDBName = "base_models_db_test"
		mongoCollectionName = "base_models_test"

		storage, connErr = NewBaseCollection(context.TODO(), mongoURI, mongoDBName, mongoCollectionName)

		if connErr != nil {
			GinkgoT().Fatal(connErr)
		}
	})

	Describe("constructor", func() {
		It("should create new storage", func() {
			Expect(storage).NotTo(BeNil())
			Expect(connErr).To(BeNil())
		})
	})

	Describe("methods", func() {
		AfterEach(func() {
			storage.DeleteAll(context.TODO())
		})

		Describe("save & find methods", func() {
			Describe(".InsertOne()", func() {
				It("should create new model", func() {
					curTime := time.Now().UTC().Add(time.Hour * 1)

					timecop.Freeze(curTime)
					defer timecop.Return()

					m := NewExampleModel()

					id, saveErr := storage.InsertOne(context.TODO(), m)

					Expect(m.CreatedAt).To(Equal(curTime))
					Expect(m.UpdatedAt).To(Equal(curTime))

					findErr := storage.GetOneByID(context.TODO(), id, m)

					Expect(id).NotTo(Equal(""))
					Expect(m.BaseDocument.GetHexID()).NotTo(Equal(""))
					Expect(saveErr).To(BeNil())
					Expect(findErr).To(BeNil())

					Expect(m.CreatedAt.Unix()).To(Equal(curTime.Unix()))
					Expect(m.UpdatedAt.Unix()).To(Equal(curTime.Unix()))
				})
			})

			Context("With existing model", func() {
				var (
					m  *ExampleModel
					id string
				)

				BeforeEach(func() {
					m = NewExampleModel()
					id, _ = storage.InsertOne(context.TODO(), m)
				})

				Describe(".GetOneByID()", func() {
					It("should find the model by id", func() {
						emptyModel := &ExampleModel{}
						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(findErr).To(BeNil())
						Expect(emptyModel.Title).To(Equal(m.Title))
					})
				})

				Describe(".UpdateOne()", func() {
					It("should update the model", func() {
						emptyModel := &ExampleModel{}

						curTime := time.Now().UTC().Add(time.Hour * 1)

						timecop.Freeze(curTime)
						defer timecop.Return()

						newTitle := "New title " + uuid.New().String()
						m.Title = newTitle

						updateErr := storage.UpdateOne(context.TODO(), m)

						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(updateErr).To(BeNil())
						Expect(findErr).To(BeNil())

						Expect(emptyModel.Title).To(Equal(newTitle))
						Expect(emptyModel.CreatedAt.Unix()).NotTo(Equal(curTime.Unix()))
						Expect(emptyModel.UpdatedAt.Unix()).To(Equal(curTime.Unix()))
					})
				})

				Describe(".ReplaceOne()", func() {
					It("should update the model", func() {
						emptyModel := &ExampleModel{}

						curTime := time.Now().UTC().Add(time.Hour * 1)

						timecop.Freeze(curTime)
						defer timecop.Return()

						newTitle := "New title " + uuid.New().String()
						m.Title = newTitle

						filter := bson.M{"_id": bson.M{"$eq": m.GetID()}}
						updateRes, updateErr := storage.ReplaceOne(context.TODO(), filter, m)
						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(updateErr).To(BeNil())
						Expect(updateRes.ModifiedCount).To(Equal(int64(1)))
						Expect(findErr).To(BeNil())

						Expect(emptyModel.Title).To(Equal(newTitle))
						Expect(emptyModel.CreatedAt.Unix()).NotTo(Equal(curTime.Unix()))
						Expect(emptyModel.UpdatedAt.Unix()).To(Equal(curTime.Unix()))
					})
				})

				Describe(".ReplaceOneByID()", func() {
					It("should update the model", func() {
						emptyModel := &ExampleModel{}

						curTime := time.Now().UTC().Add(time.Hour * 1)

						timecop.Freeze(curTime)
						defer timecop.Return()

						newTitle := "New title " + uuid.New().String()
						m.Title = newTitle

						updateRes, updateErr := storage.ReplaceOneByID(context.TODO(), m.GetHexID(), m)
						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(updateErr).To(BeNil())
						Expect(updateRes.ModifiedCount).To(Equal(int64(1)))
						Expect(findErr).To(BeNil())
						Expect(emptyModel.Title).To(Equal(newTitle))
						Expect(emptyModel.CreatedAt.Unix()).NotTo(Equal(curTime.Unix()))
						Expect(emptyModel.UpdatedAt.Unix()).To(Equal(curTime.Unix()))
					})
				})
			})

			Context("without any stored models", func() {
				var (
					m  *ExampleModel
					id string
				)

				BeforeEach(func() {
					id = "123"
					m = NewExampleModel()
				})

				Describe(".GetOneByID()", func() {
					It("should find the model by id", func() {
						emptyModel := &ExampleModel{}
						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(findErr).To(Equal(ErrInvalidObjectID))
					})
				})

				Describe(".UpdateOne()", func() {
					It("should update the model", func() {
						emptyModel := &ExampleModel{}
						updateErr := storage.UpdateOne(context.TODO(), m)
						findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

						Expect(updateErr).NotTo(BeNil())
						Expect(findErr).To(Equal(ErrInvalidObjectID))
					})
				})
			})
		})

		Describe(".InsertMany()", func() {
			It("should insert many records", func() {
				curTime := time.Now().UTC().Add(time.Hour * 1)

				timecop.Freeze(curTime)
				defer timecop.Return()

				docs := []interface{}{
					NewExampleModel(),
					NewExampleModel(),
					NewExampleModel(),
				}

				for i, _ := range docs {
					docs[i].(*ExampleModel).SetupCreatedAt()
					docs[i].(*ExampleModel).SetupUpdatedAt()
				}

				ids, err := storage.InsertMany(context.TODO(), docs)

				Expect(err).To(BeNil())
				Expect(len(ids)).To(Equal(len(docs)))

				for _, id := range ids {
					emptyModel := &ExampleModel{}
					findErr := storage.GetOneByID(context.TODO(), id, emptyModel)

					Expect(findErr).To(BeNil())

					Expect(emptyModel.CreatedAt.Unix()).To(Equal(curTime.Unix()))
					Expect(emptyModel.UpdatedAt.Unix()).To(Equal(curTime.Unix()))
				}
			})
		})

		Describe("UpdateMany()", func() {
			var (
				ids       []string
				objectIDs []primitive.ObjectID
				docs      []interface{}
				curTime   time.Time
			)

			BeforeEach(func() {
				curTime = time.Now().UTC().Add(time.Hour * 1)
				docs = []interface{}{
					NewExampleModel(),
					NewExampleModel(),
					NewExampleModel(),
				}

				timecop.Freeze(curTime)
				for i, _ := range docs {
					docs[i].(*ExampleModel).SetupCreatedAt()
					docs[i].(*ExampleModel).SetupUpdatedAt()
				}
				timecop.Return()

				ids, _ = storage.InsertMany(context.TODO(), docs)
			})

			It("should update many records", func() {
				timecop.Freeze(curTime)
				defer timecop.Return()

				newTime := curTime.Add(time.Second * 10)

				for _, id := range ids {
					objID, _ := primitive.ObjectIDFromHex(id)
					objectIDs = append(objectIDs, objID)
				}

				filter := bson.M{"_id": bson.M{"$in": objectIDs}}
				update := bson.M{"$set": bson.M{"updated_at": newTime}}

				res, err := storage.UpdateMany(
					context.TODO(),
					filter,
					update,
				)

				Expect(err).To(BeNil())

				docsCount := int64(len(docs))

				Expect(res.MatchedCount).To(Equal(docsCount))
				Expect(res.ModifiedCount).To(Equal(docsCount))

				l, err := storage.GetManyByFilter(
					context.TODO(),
					filter,
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(err).To(BeNil())

				for _, m := range l {
					Expect(m.(*ExampleModel).UpdatedAt.Unix()).To(Equal(newTime.Unix()))
				}
			})
		})

		Describe(".FindAllByFilter()", func() {
			var (
				m1, m2, m3 *ExampleModel
				l          []*ExampleModel
			)

			BeforeEach(func() {
				m1 = NewExampleModel()
				m2 = NewExampleModel()
				m3 = NewExampleModel()

				storage.InsertOne(context.TODO(), m1)
				storage.InsertOne(context.TODO(), m2)
				storage.InsertOne(context.TODO(), m3)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should return all models", func() {
				err := storage.FindAllByFilter(
					context.TODO(),
					bson.M{},
					&l,
				)

				Expect(err).To(BeNil())
				Expect(l).NotTo(BeNil())
				Expect(len(l)).To(Equal(3))
			})
		})

		Describe(".GetManyByFilter()", func() {
			var (
				m1, m2, m3 *ExampleModel
			)

			BeforeEach(func() {
				m1 = NewExampleModel()
				m2 = NewExampleModel()
				m3 = NewExampleModel()

				storage.InsertOne(context.TODO(), m1)
				storage.InsertOne(context.TODO(), m2)
				storage.InsertOne(context.TODO(), m3)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should return all models", func() {
				l, err := storage.GetManyByFilter(
					context.TODO(),
					bson.M{},
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(err).To(BeNil())
				Expect(l).NotTo(BeNil())
				Expect(len(l)).To(Equal(3))
			})
		})

		Describe(".CountByFilter()", func() {
			var (
				m1, m2, m3 *ExampleModel
			)

			BeforeEach(func() {
				m1 = NewExampleModel()
				m2 = NewExampleModel()
				m3 = NewExampleModel()

				storage.InsertOne(context.TODO(), m1)
				storage.InsertOne(context.TODO(), m2)
				storage.InsertOne(context.TODO(), m3)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should return count of documents", func() {
				count, err := storage.CountByFilter(
					context.TODO(),
					bson.M{},
				)

				Expect(err).To(BeNil())
				Expect(count).To(Equal(int64(3)))
			})
		})

		Describe(".DeleteManyByFilter()", func() {
			var (
				m1, m2, m3 *ExampleModel
			)

			BeforeEach(func() {
				m1 = NewExampleModel()
				m2 = NewExampleModel()
				m3 = NewExampleModel()

				storage.InsertOne(context.TODO(), m1)
				storage.InsertOne(context.TODO(), m2)
				storage.InsertOne(context.TODO(), m3)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should delete model by filter", func() {
				l1, err1 := storage.GetManyByFilter(
					context.TODO(),
					bson.M{},
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(err1).To(BeNil())
				Expect(len(l1)).To(Equal(3))

				res, delErr := storage.DeleteManyByFilter(context.TODO(), bson.M{
					"title": m1.Title,
				})

				Expect(delErr).To(BeNil())

				l2, err2 := storage.GetManyByFilter(
					context.TODO(),
					bson.M{},
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(l2).NotTo(BeNil())
				Expect(err2).To(BeNil())
				Expect(len(l2)).To(Equal(2))
				Expect(res.DeletedCount).To(Equal(int64(1)))
			})
		})

		Describe(".DeleteOneByID()", func() {
			var (
				m1, m2, m3 *ExampleModel
			)

			BeforeEach(func() {
				m1 = NewExampleModel()
				m2 = NewExampleModel()
				m3 = NewExampleModel()

				storage.InsertOne(context.TODO(), m1)
				storage.InsertOne(context.TODO(), m2)
				storage.InsertOne(context.TODO(), m3)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should delete model by id", func() {
				l1, err1 := storage.GetManyByFilter(
					context.TODO(),
					bson.M{},
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(err1).To(BeNil())
				Expect(len(l1)).To(Equal(3))

				delErr := storage.DeleteOneByID(context.TODO(), m1.GetHexID())

				Expect(delErr).To(BeNil())

				l2, err2 := storage.GetManyByFilter(
					context.TODO(),
					bson.M{},
					func() Document {
						return &ExampleModel{}
					},
				)

				Expect(l2).NotTo(BeNil())
				Expect(err2).To(BeNil())
				Expect(len(l2)).To(Equal(2))
			})
		})

		Describe(".UpsertOne()", func() {
			var (
				m1    *ExampleModel
				docID string
			)

			BeforeEach(func() {
				docID = "555555555555555555555555"
				m1 = NewExampleModel()
				m1.SetHexID(docID)
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should upsert model if was no previous", func() {
				ctx := context.Background()
				filter := bson.M{"_id": m1.BaseDocument.ID}
				update := bson.M{
					"$set": bson.M{
						"title": m1.Title,
					},
				}

				/* insert record */
				insertedModel, err := storage.UpsertOne(ctx, filter, update, &ExampleModel{})
				Expect(err).To(BeNil())
				Expect(m1).To(Equal(insertedModel))

				findRes := &ExampleModel{}
				findErr := storage.GetOneByID(context.TODO(), docID, findRes)
				Expect(findErr).To(BeNil())
				Expect(findRes).To(Equal(insertedModel))

				update = bson.M{
					"$set": bson.M{
						"title": m1.Title + "updated",
					},
				}

				/* upsert existing record */
				upsertedModel, err := storage.UpsertOne(ctx, filter, update, &ExampleModel{})
				Expect(err).To(BeNil())

				m1.Title = m1.Title + "updated"
				Expect(upsertedModel).To(Equal(m1))

				findErr = storage.GetOneByID(context.TODO(), docID, findRes)
				Expect(findErr).To(BeNil())
				Expect(findRes).To(Equal(upsertedModel))
			})
		})

		Describe(".FindAndUpdateOne()", func() {
			var (
				m1    *ExampleModel
				docID string
			)

			BeforeEach(func() {
				m1 = NewExampleModel()

				insertedID, err := storage.InsertOne(context.TODO(), m1)
				Expect(err).To(BeNil())

				docID = insertedID
			})

			AfterEach(func() {
				storage.DeleteAll(context.TODO())
			})

			It("should update model when it exists", func() {
				ctx := context.Background()
				filter := bson.M{"_id": m1.BaseDocument.ID}
				update := bson.M{
					"$set": bson.M{
						"title": m1.Title,
					},
				}

				update = bson.M{
					"$set": bson.M{
						"title": m1.Title + "updated",
					},
				}

				/* update existing record */
				updatedModel, err := storage.FindAndUpdateOne(ctx, filter, update, &ExampleModel{})
				Expect(err).To(BeNil())

				Expect(updatedModel.(*ExampleModel).Title).To(Equal(m1.Title + "updated"))

				findRes := &ExampleModel{}
				findErr := storage.GetOneByID(context.TODO(), docID, findRes)
				Expect(findErr).To(BeNil())
				Expect(findRes).To(Equal(updatedModel))
			})
		})
	})
})

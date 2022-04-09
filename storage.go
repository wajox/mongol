package mongol

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage
type Storage interface {
	AddBeforeHook(methodName string, h Hook)
	AddAfterHook(methodName string, h Hook)
	Ping(ctx context.Context) error
	Collection() *mongo.Collection
	Database() *mongo.Database
	MongoClient() *mongo.Client
	CreateIndex(ctx context.Context, k interface{}, o *options.IndexOptions) (string, error)
	InsertOne(ctx context.Context, m Document, opts ...*options.InsertOneOptions) (string, error)
	InsertMany(ctx context.Context, docs []interface{}, opts ...*options.InsertManyOptions) ([]string, error)
	UpdateOne(ctx context.Context, m Document, opts ...*options.UpdateOptions) error
	UpdateManyByFilter(ctx context.Context, filter interface{}, m Document, opts ...*options.UpdateOptions) error
	UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpsertOne(ctx context.Context, filter interface{}, update bson.M, m Document) (Document, error)
	FindAndUpdateOne(ctx context.Context, filter interface{}, update bson.M, m Document) (Document, error)
	ReplaceOne(ctx context.Context, filter interface{}, m Document, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	ReplaceOneByID(ctx context.Context, recordID string, m Document, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	GetOneByID(ctx context.Context, recordID string, m Document, opts ...*options.FindOneOptions) error
	GetOneByFilter(ctx context.Context, filter interface{}, m Document, opts ...*options.FindOneOptions) error
	GetManyByFilter(ctx context.Context, filter interface{}, modelBuilder func() Document, opts ...*options.FindOptions) ([]Document, error)
	FindAllByFilter(ctx context.Context, filter interface{}, docs interface{}, opts ...*options.FindOptions) error
	FindManyByFilter(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	CountByFilter(ctx context.Context, filter interface{}) (int64, error)
	DeleteManyByFilter(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOneByID(ctx context.Context, docID string) error
	DeleteAll(ctx context.Context) error
}

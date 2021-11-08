package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionIDKey          = "_id"
	CreateIndexMethod        = "CreateIndex"
	InsertOneMethod          = "InsertOne"
	InsertManyMethod         = "InsertMany"
	UpdateOneMethod          = "UpdateOne"
	UpdateManyMethod         = "UpdateMany"
	UpdateManyByFilterMethod = "UpdateManyByFilter"
	ReplaceOneMethod         = "ReplaceOne"
	ReplaceOneByIDMethod     = "ReplaceOneByID"
	GetOneByIDMethod         = "GetOneByID"
	GetOneByFilterMethod     = "GetOneByFilter"
	GetManyByFilterMethod    = "GetManyByFilter"
	FindAllByFilterMethod    = "FindAllByFilter"
	FindManyByFilterMethod   = "FindManyByFilter"
	UpsertOneMethod          = "UpsertOne"
	DeleteManyMethod         = "DeleteMany"
	DeleteManyByFilterMethod = "DeleteManyByFilter"
	DeleteOneByIDMethod      = "DeleteOneByID"
	DeleteAllMethod          = "DeleteAll"
	CloseCursorTimeout       = time.Second * 1
	FetchTimeout             = time.Second * 1
	QueryTimeout             = time.Second * 1
	FilterTimeout            = time.Second * 1
)

type Hook func(ctx context.Context) error

type Client struct {
	mongoClient *mongo.Client
}

func (c *Client) GetMongoClient() *mongo.Client {
	return c.mongoClient
}

// BaseCollectionStorage
type BaseCollectionStorage struct {
	Client         *Client
	DBName         string
	CollectionName string
	BeforeHooks    map[string][]Hook
	AfterHooks     map[string][]Hook
}

// Document
type Document interface {
	GetID() primitive.ObjectID
	GetHexID() string
	SetHexID(hexID string) error
	SetJSONID(jsonB []byte) error
	SetupTimestamps()
}

func NewClient(ctx context.Context, mongoURI string) (*Client, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	return &Client{mongoClient: c}, nil
}

// NewBaseCollectionStorage() is a constructor for BaseCollectionStorage struct
func NewBaseCollectionStorage(ctx context.Context, mongoURI, dbName, collectionName string) (*BaseCollectionStorage, error) {
	client, err := NewClient(ctx, mongoURI)
	if err != nil {
		return nil, err
	}

	return NewBaseCollectionStorageWithClient(client, dbName, collectionName), nil
}

func NewBaseCollectionStorageWithClient(client *Client, dbName, collectionName string) *BaseCollectionStorage {
	return &BaseCollectionStorage{
		Client:         client,
		DBName:         dbName,
		CollectionName: collectionName,
		BeforeHooks:    make(map[string][]Hook),
		AfterHooks:     make(map[string][]Hook),
	}
}

func (s *BaseCollectionStorage) runBeforeHooks(ctx context.Context, methodName string) error {
	hooks, ok := s.BeforeHooks[methodName]
	if !ok {
		return nil
	}

	for i := range hooks {
		if err := hooks[i](ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *BaseCollectionStorage) runAfterHooks(ctx context.Context, methodName string) error {
	hooks, ok := s.AfterHooks[methodName]
	if !ok {
		return nil
	}

	for i := range hooks {
		if err := hooks[i](ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *BaseCollectionStorage) AddBeforeHook(methodName string, h Hook) {
	if _, ok := s.BeforeHooks[methodName]; !ok {
		s.BeforeHooks[methodName] = []Hook{}
	}

	s.BeforeHooks[methodName] = append(s.BeforeHooks[methodName], h)
}

func (s *BaseCollectionStorage) AddAfterHook(methodName string, h Hook) {
	if _, ok := s.AfterHooks[methodName]; !ok {
		s.AfterHooks[methodName] = []Hook{}
	}

	s.AfterHooks[methodName] = append(s.AfterHooks[methodName], h)
}

// Ping() the mongo server
func (s *BaseCollectionStorage) Ping(ctx context.Context) error {
	return s.Client.GetMongoClient().Ping(ctx, nil)
}

// GetCollection() returns storage collection
func (s *BaseCollectionStorage) GetCollection() *mongo.Collection {
	return s.Client.GetMongoClient().Database(s.DBName).Collection(s.CollectionName)
}

func (s *BaseCollectionStorage) CreateIndex(ctx context.Context, k interface{}, o *options.IndexOptions) (string, error) {
	if err := s.runBeforeHooks(ctx, CreateIndexMethod); err != nil {
		return "", err
	}

	defer s.runAfterHooks(ctx, CreateIndexMethod)

	return s.GetCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    k,
		Options: o,
	})
}

// InsertOne() inserts given Document and returns an ID of inserted document
func (s *BaseCollectionStorage) InsertOne(ctx context.Context, m Document, opts ...*options.InsertOneOptions) (string, error) {
	if err := s.runBeforeHooks(ctx, InsertOneMethod); err != nil {
		return "", err
	}
	defer s.runAfterHooks(ctx, InsertOneMethod)

	m.SetupTimestamps()

	b, err := bson.Marshal(m)
	if err != nil {
		return "", err
	}

	res, err := s.GetCollection().InsertOne(ctx, b, opts...)
	if err != nil {
		return "", HandleDuplicationErr(err)
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", ErrInvalidObjectID
	}

	hexID := objectID.Hex()
	if err := m.SetHexID(hexID); err != nil {
		return "", err
	}

	return hexID, nil
}

// InsertMany()
func (s *BaseCollectionStorage) InsertMany(ctx context.Context, docs []interface{}, opts ...*options.InsertManyOptions) ([]string, error) {
	if err := s.runBeforeHooks(ctx, InsertManyMethod); err != nil {
		return []string{}, err
	}
	defer s.runAfterHooks(ctx, InsertManyMethod)

	res, err := s.GetCollection().InsertMany(ctx, docs, opts...)

	if err != nil {
		return []string{}, HandleDuplicationErr(err)
	}

	hexIDs := make([]string, len(res.InsertedIDs))

	for i, insertedID := range res.InsertedIDs {
		objectID, ok := insertedID.(primitive.ObjectID)

		if !ok {
			return hexIDs, ErrInvalidObjectID
		}

		hexIDs[i] = objectID.Hex()
	}

	return hexIDs, nil
}

// UpdateOne() updates given Document
func (s *BaseCollectionStorage) UpdateOne(ctx context.Context, m Document, opts ...*options.UpdateOptions) error {
	if err := s.runBeforeHooks(ctx, UpdateOneMethod); err != nil {
		return err
	}
	defer s.runAfterHooks(ctx, UpdateOneMethod)

	filter := bson.M{CollectionIDKey: bson.M{"$eq": m.GetID()}}
	return s.UpdateManyByFilter(ctx, filter, m, opts...)
}

// UpdateByFilter() updates given Document according to provided filter
func (s *BaseCollectionStorage) UpdateManyByFilter(ctx context.Context, filter interface{}, m Document, opts ...*options.UpdateOptions) error {
	if err := s.runBeforeHooks(ctx, UpdateManyByFilterMethod); err != nil {
		return err
	}
	defer s.runAfterHooks(ctx, UpdateManyByFilterMethod)

	m.SetupTimestamps()

	res, err := s.GetCollection().UpdateOne(
		ctx,
		filter,
		bson.D{primitive.E{Key: "$set", Value: m}},
		opts...,
	)
	if err != nil {
		return HandleDuplicationErr(err)
	}

	if res.MatchedCount != 1 {
		return ErrDocumentNotFound
	}

	if res.ModifiedCount != 1 {
		return ErrDocumentNotModified
	}

	return nil
}

// UpdateMany()
func (s *BaseCollectionStorage) UpdateMany(
	ctx context.Context,
	filter, update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	if err := s.runBeforeHooks(ctx, UpdateManyMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, UpdateManyMethod)

	return s.GetCollection().UpdateMany(
		ctx,
		filter,
		update,
		opts...,
	)
}

// UpsertOne - insert or update existing record. Returns updated model.
func (s *BaseCollectionStorage) UpsertOne(
	ctx context.Context,
	filter interface{},
	update bson.M,
	m Document,
) (Document, error) {
	if err := s.runBeforeHooks(ctx, UpsertOneMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, UpsertOneMethod)

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(true)

	res := s.GetCollection().FindOneAndUpdate(ctx, filter, update, opts)
	if err := res.Decode(m); err != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, ErrDocumentNotFound
		}

		return nil, err
	}

	b, err := res.DecodeBytes()
	if err != nil {
		return nil, err
	}

	if err := m.SetJSONID(b.Lookup(CollectionIDKey).Value); err != nil {
		return nil, err
	}

	return m, nil
}

// ReplaceOne
func (s *BaseCollectionStorage) ReplaceOne(
	ctx context.Context,
	filter interface{},
	m Document,
	opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {
	if err := s.runBeforeHooks(ctx, ReplaceOneMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, ReplaceOneMethod)

	m.SetupTimestamps()

	return s.GetCollection().ReplaceOne(
		ctx,
		filter,
		m,
		opts...,
	)
}

// ReplaceOneByID()
func (s *BaseCollectionStorage) ReplaceOneByID(
	ctx context.Context,
	recordID string,
	m Document,
	opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {
	if err := s.runBeforeHooks(ctx, ReplaceOneByIDMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, ReplaceOneByIDMethod)

	oid, err := primitive.ObjectIDFromHex(recordID)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	return s.ReplaceOne(
		ctx,
		bson.M{CollectionIDKey: bson.M{"$eq": oid}},
		m,
		opts...,
	)
}

// GetOneByID() is trying to find Document by given recordID
func (s *BaseCollectionStorage) GetOneByID(
	ctx context.Context,
	recordID string,
	m Document,
	opts ...*options.FindOneOptions,
) error {
	if err := s.runBeforeHooks(ctx, GetOneByIDMethod); err != nil {
		return err
	}

	defer s.runAfterHooks(ctx, GetOneByIDMethod)

	oid, err := primitive.ObjectIDFromHex(recordID)
	if err != nil {
		return ErrInvalidObjectID
	}

	filter := bson.M{CollectionIDKey: oid}

	return s.GetOneByFilter(ctx, filter, m, opts...)
}

// GetOneByFilter() is trying to find Document by provided filter
func (s *BaseCollectionStorage) GetOneByFilter(
	ctx context.Context,
	filter interface{},
	m Document,
	opts ...*options.FindOneOptions,
) error {
	if err := s.runBeforeHooks(ctx, GetOneByFilterMethod); err != nil {
		return err
	}

	defer s.runAfterHooks(ctx, GetOneByFilterMethod)

	res := s.GetCollection().FindOne(ctx, filter, opts...)
	if err := res.Decode(m); err != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return ErrDocumentNotFound
		}

		return err
	}

	b, err := res.DecodeBytes()

	if err != nil {
		return err
	}

	if err := m.SetJSONID(b.Lookup(CollectionIDKey).Value); err != nil {
		return err
	}

	return nil
}

// GetManyByFilter()
func (s *BaseCollectionStorage) GetManyByFilter(
	ctx context.Context,
	filter interface{},
	modelBuilder func() Document,
	opts ...*options.FindOptions,
) ([]Document, error) {
	if err := s.runBeforeHooks(ctx, GetManyByFilterMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, GetManyByFilterMethod)

	filterCtx, filterCancel := context.WithTimeout(ctx, FilterTimeout)
	defer filterCancel()

	cur, err := s.FindManyByFilter(filterCtx, filter, opts...)
	if err != nil {
		return nil, err
	}

	closeCtx, closeCancel := context.WithTimeout(ctx, CloseCursorTimeout)
	defer closeCancel()
	defer cur.Close(closeCtx)

	var l []Document

	nextCtx, nextCancel := context.WithTimeout(context.Background(), FetchTimeout)

	defer nextCancel()

	for cur.Next(nextCtx) {
		m := modelBuilder()
		if err := cur.Decode(m); err != nil {
			return nil, err
		}

		if err := m.SetJSONID(cur.Current.Lookup(CollectionIDKey).Value); err != nil {
			return nil, err
		}

		l = append(l, m)
	}

	return l, nil
}

// FindAllByFilter()
func (s *BaseCollectionStorage) FindAllByFilter(
	ctx context.Context,
	filter interface{},
	docs interface{},
	opts ...*options.FindOptions,
) error {
	if err := s.runBeforeHooks(ctx, FindAllByFilterMethod); err != nil {
		return err
	}

	defer s.runAfterHooks(ctx, FindAllByFilterMethod)

	filterCtx, filterCancel := context.WithTimeout(ctx, FilterTimeout)
	defer filterCancel()

	cur, err := s.FindManyByFilter(filterCtx, filter, opts...)
	if err != nil {
		return err
	}

	closeCtx, closeCancel := context.WithTimeout(ctx, CloseCursorTimeout)
	defer closeCancel()
	defer cur.Close(closeCtx)

	allCtx, allCancel := context.WithTimeout(context.Background(), FetchTimeout)
	defer allCancel()

	return cur.All(allCtx, docs)
}

// FindManyByFilter()
func (s *BaseCollectionStorage) FindManyByFilter(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (*mongo.Cursor, error) {
	if err := s.runBeforeHooks(ctx, FindManyByFilterMethod); err != nil {
		return nil, err
	}

	defer s.runAfterHooks(ctx, FindManyByFilterMethod)

	cur, err := s.GetCollection().Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	if cur.Err() == nil {
		return cur, nil
	}

	closeCtx, closeCancel := context.WithTimeout(ctx, CloseCursorTimeout)
	defer closeCancel()
	defer cur.Close(closeCtx)

	return nil, cur.Err()
}

// DeleteManyByFilter() documents by given filters
func (s *BaseCollectionStorage) DeleteManyByFilter(
	ctx context.Context,
	filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	if err := s.runBeforeHooks(ctx, DeleteManyByFilterMethod); err != nil {
		return nil, err
	}
	defer s.runAfterHooks(ctx, DeleteManyByFilterMethod)

	return s.GetCollection().DeleteMany(ctx, filter, opts...)
}

// DeleteOneByID() deletes document by given ID
func (s *BaseCollectionStorage) DeleteOneByID(ctx context.Context, docID string) error {
	if err := s.runBeforeHooks(ctx, DeleteOneByIDMethod); err != nil {
		return err
	}

	defer s.runAfterHooks(ctx, DeleteOneByIDMethod)

	oid, err := primitive.ObjectIDFromHex(docID)
	if err != nil {
		return ErrInvalidObjectID
	}

	filter := bson.M{CollectionIDKey: oid}

	r, err := s.DeleteManyByFilter(ctx, filter)
	if r.DeletedCount != 1 {
		return ErrDocumentNotFound
	}

	return err
}

// DropAll() deletes collection from database
func (s *BaseCollectionStorage) DeleteAll(ctx context.Context) error {
	if err := s.runBeforeHooks(ctx, DeleteAllMethod); err != nil {
		return err
	}

	defer s.runAfterHooks(ctx, DeleteAllMethod)

	return s.GetCollection().Drop(ctx)
}

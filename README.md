# Mongol

[![Go Reference](https://pkg.go.dev/badge/github.com/wajox/mongol.svg)](https://pkg.go.dev/github.com/wajox/mongol)

## Model Example
```golang
type ExampleModel struct {
	mongol.BaseDocument `bson:",inline"`
	// Default fields from BaseDocument
	//
	// ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	// UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	Title string `json:"title,omitempty" bson:"title,omitempty"`
}
```

## Storage usage example
```golang
storage, connErr = NewBaseCollection(context.TODO(), mongoURI, mongoDBName, mongoCollectionName)

m := &ExampleModel{}

id, saveErr := storage.InsertOne(context.TODO(), m)
```

# Run tests

```
make test-all
```

```
MONGODB_URI=mongodb://localhost:27017 go test -a -v ./...
```
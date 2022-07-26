# Mongol

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/wajox/mongol.svg)](https://pkg.go.dev/github.com/wajox/mongol)
[![codecov](https://codecov.io/gh/wajox/mongol/branch/master/graph/badge.svg?token=MFEF13319U)](https://codecov.io/gh/wajox/mongol)
![Go workflow](https://github.com/wajox/mongol/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wajox/mongol)](https://goreportcard.com/report/github.com/wajox/mongol)

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

## Fiilter example
```golang

newM := &ExampleModel{}

err := storage.GetOneByFilter(
	context.TODO(),
	NewFilterBuilder().EqualTo("_id", m.ID).GetQuery(),
	newM,
)

if err != nil {
	fmt.Printf("%#v", newM)
}
```

# Run tests

```
make test-all
```

```
MONGODB_URI=mongodb://localhost:27017 go test -a -v ./...
```
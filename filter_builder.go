package mongol

import (
	"go.mongodb.org/mongo-driver/bson"
)

// FilterBuilder
type FilterBuilder struct {
	query bson.M
}

// NewFilterBuilder() initializes a new FilterBuilder
func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		query: bson.M{},
	}
}

// GetQuery() returns compiled query for *FilterBuilder
func (fb *FilterBuilder) GetQuery() interface{} {
	return fb.query
}

// Where() allows to set custom conditions for the query
func (fb *FilterBuilder) Where(key string, query bson.M) *FilterBuilder {
	fb.query[key] = query

	return fb
}

// Or() allows to combine few queries with $or
func (fb *FilterBuilder) Or(query ...interface{}) *FilterBuilder {
	if _, exists := fb.query["$or"]; !exists {
		fb.query["$or"] = bson.A{}
	}

	for i := range query {
		fb.query["$or"] = append(fb.query["$or"].(bson.A), query[i])
	}

	return fb
}

// And() allows to combine few queries with $and
func (fb *FilterBuilder) And(query ...interface{}) *FilterBuilder {
	if _, exists := fb.query["$and"]; !exists {
		fb.query["$and"] = bson.A{}
	}

	for i := range query {
		fb.query["$and"] = append(fb.query["$and"].(bson.A), query[i])
	}

	return fb
}

// EqualTo() implements $eq condition for the query
func (fb *FilterBuilder) EqualTo(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$eq": value}

	return fb
}

// NotEqualTo() implements $ne condition for the query
func (fb *FilterBuilder) NotEqualTo(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$ne": value}

	return fb
}

// In() implements $in condition for the query
func (fb *FilterBuilder) In(key string, values []interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$in": values}

	return fb
}

// NotIn() implements $nin condition for the query
func (fb *FilterBuilder) NotIn(key string, values []interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$nin": values}

	return fb
}

// HasField() implements $exists:true condition for the query
func (fb *FilterBuilder) HasField(key string) *FilterBuilder {
	fb.query[key] = bson.M{"$exists": true}

	return fb
}

// HasNotField() implements $exists:false condition for the query
func (fb *FilterBuilder) HasNotField(key string) *FilterBuilder {
	fb.query[key] = bson.M{"$exists": false}

	return fb
}

// Gte() implements $gte condition for the query
func (fb *FilterBuilder) Gte(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$gte": value}

	return fb
}

// Lte() implements $lte condition for the query
func (fb *FilterBuilder) Lte(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$lte": value}

	return fb
}

// Gt() implements $gt condition for the query
func (fb *FilterBuilder) Gt(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$gt": value}

	return fb
}

// Lt() implements $lt condition for the query
func (fb *FilterBuilder) Lt(key string, value interface{}) *FilterBuilder {
	fb.query[key] = bson.M{"$lt": value}

	return fb
}

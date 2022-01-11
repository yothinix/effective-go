package store

import (
	"apidemo/todo"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStore struct{
	*mongo.Collection
}

func NewMongoDBStore(col *mongo.Collection) *MongoDBStore {
	return &MongoDBStore{Collection: col}
}

func (s *MongoDBStore) New(todo *todo.Todo) error {
	_, err := s.Collection.InsertOne(context.Background(), todo)
	return err
}

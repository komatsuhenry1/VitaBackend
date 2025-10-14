package repository

import (
	"context"
	"medassist/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	FindMessagesBetween(userID, otherUserID primitive.ObjectID) ([]model.Message, error)
}

type messageRepositoryImpl struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepositoryImpl{
		collection: db.Collection("messages"),
	}
}

func (r *messageRepositoryImpl) FindMessagesBetween(userID, otherUserID primitive.ObjectID) ([]model.Message, error) {
	var messages []model.Message
	ctx := context.TODO() 

	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID, "receiver_id": otherUserID},
			{"sender_id": otherUserID, "receiver_id": userID},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) 

	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
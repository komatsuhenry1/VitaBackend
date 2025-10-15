package repository

import (
	"context"
	"medassist/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"medassist/internal/chat/dto"
)

type MessageRepository interface {
	FindMessagesBetween(userID, otherUserID primitive.ObjectID) ([]model.Message, error)
	Save(message *model.Message) error
	GetConversationsForUser(userID primitive.ObjectID) ([]dto.ConversationDTO, error)
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

func (r *messageRepositoryImpl) Save(message *model.Message) error {
	ctx := context.TODO()
	// Define o ID e o Timestamp antes de inserir
	message.ID = primitive.NewObjectID()
	message.Timestamp = time.Now()

	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *messageRepositoryImpl) GetConversationsForUser(userID primitive.ObjectID) ([]dto.ConversationDTO, error) {
	ctx := context.TODO()
	var conversations []dto.ConversationDTO

	// O Aggregation Pipeline é uma sequência de etapas para transformar os dados
	pipeline := mongo.Pipeline{
		// Etapa 1: Encontra todas as mensagens onde o usuário é remetente OU destinatário
		{{"$match", bson.M{
			"$or": []bson.M{
				{"sender_id": userID},
				{"receiver_id": userID},
			},
		}}},
		// Etapa 2: Ordena as mensagens da mais nova para a mais antiga
		{{"$sort", bson.M{"timestamp": -1}}},
		// Etapa 3: Agrupa as mensagens por conversa (pela outra pessoa)
		{{"$group", bson.M{
			// Identificador do grupo: quem é a outra pessoa na conversa
			"_id": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": []interface{}{"$sender_id", userID}},
					"then": "$receiver_id",
					"else": "$sender_id",
				},
			},
			// Pega a primeira mensagem do grupo (que é a mais recente, devido à ordenação)
			"last_message":           bson.M{"$first": "$content"},
			"last_message_timestamp": bson.M{"$first": "$timestamp"},
		}}},
		// Etapa 4: Faz um "join" com a coleção de usuários para pegar o nome e a imagem do parceiro de chat
		{{"$lookup", bson.M{
			"from":         "users", // O nome da sua coleção de usuários
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "partnerInfo",
		}}},
		// Etapa 5: Desconstrói o array 'partnerInfo' para que possamos acessar seus campos
		{{"$unwind", "$partnerInfo"}},
		// Etapa 6: Projeta (formata) o resultado final para corresponder ao nosso DTO
		{{"$project", bson.M{
			"_id":                    0,
			"partner_id":             "$_id",
			"partner_name":           "$partnerInfo.name",
			"partner_image_id":       "$partnerInfo.profile_image_id",
			"last_message":           "$last_message",
			"last_message_timestamp": "$last_message_timestamp",
		}}},
		// Etapa 7: Ordena as conversas pela mais recente
		{{"$sort", bson.M{"last_message_timestamp": -1}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}


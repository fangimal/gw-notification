package mongodb

import (
	"context"
	"gw-notification/pkg/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewStorage(uri string) (*Storage, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Storage{
		client: client,
		db:     client.Database("notifications"),
	}, nil
}

func (s *Storage) SaveNotification(ctx context.Context, n *models.Notification) error {
	_, err := s.db.Collection("transfers").InsertOne(ctx, n)
	return err
}

func (s *Storage) Close() error {
	return s.client.Disconnect(context.TODO())
}
